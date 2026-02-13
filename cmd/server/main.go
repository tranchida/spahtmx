package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"spahtmx/internal/adapter/database"
	"spahtmx/internal/adapter/web"
	"spahtmx/internal/app"
	"spahtmx/internal/config"
	"spahtmx/internal/domain"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db := initDB(ctx, cfg)
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Error disconnecting from PostgreSQL", "error", err)
		}
	}()

	userRepo := database.UserBunRepository{
		DB: db,
	}

	userService := app.NewUserService(userRepo)

	prizeRepo := &database.PrizeBunRepository{
		DB: db,
	}

	prizeService := app.NewPrizeService(prizeRepo)
	authService := app.NewAuthService(userRepo)

	e := initWeb(userService, prizeService, authService)

	// Démarrage du serveur dans une goroutine
	go func() {
		slog.Info("Server starting", "url", "http://localhost:"+cfg.Port)
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			slog.Error("Server start failed", "error", err)
			os.Exit(1)
		}
	}()

	// Attente du signal d'arrêt
	<-ctx.Done()
	slog.Info("Shutting down server...")

	// Arrêt gracieux du serveur Web avec un timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}

func initDB(ctx context.Context, cfg *config.Config) *bun.DB {

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(cfg.DatabaseURL),
	))
	db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())

	if cfg.DebugSQL {
		db = db.WithQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
		))
	}

	if err := createSchema(ctx, db); err != nil {
		slog.Error("Failed to create schema", "error", err)
		os.Exit(1)
	}

	if cfg.SeedDB {
		seedUserDatabase(ctx, db)
		seedPrizeDatabase(ctx, db)
	}

	return db
}

func createSchema(ctx context.Context, db *bun.DB) error {
	models := []interface{}{
		(*database.UserBun)(nil),
		(*database.PrizeBun)(nil),
		(*database.LaureateBun)(nil),
	}

	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedUserDatabase(ctx context.Context, db *bun.DB) {
	// Mot de passe "password" haché : $2a$10$Un8S9v2vDqT5v.vQJ2vOLeC9L/9e6Z/v9.v/v9.v/v9.v/v9.v/v9
	// On va utiliser bcrypt pour générer un vrai hash pour "password"
	password := "password"
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	hashedPassword := string(bytes)

	us := []database.UserBun{
		{Username: "alice", Password: hashedPassword, Email: "alice@fake.com", Status: true},
		{Username: "bob", Password: hashedPassword, Email: "bob@fake.com", Status: false},
		{Username: "charlie", Password: hashedPassword, Email: "charlie@fake.com", Status: true},
	}

	for _, u := range us {
		_, err := db.NewInsert().Model(&u).Exec(ctx)
		if err != nil {
			slog.Error("Failed to insert user", "error", err)
			os.Exit(1)
		}
	}

	slog.Info("Database seeded successfully")
}

func seedPrizeDatabase(ctx context.Context, db *bun.DB) {

	data, err := os.ReadFile("nobel-prize.json")
	if err != nil {
		slog.Error("failed to read novel-prize.json", "error", err)
	}

	var pl domain.PrizeList
	if err := json.Unmarshal(data, &pl); err != nil {
		slog.Error("failed to unmarshal JSON", "error", err)
	}

	prizes := pl.Prizes
	fmt.Printf("Loaded %d prizes\n", len(prizes))

	prizeRepo := database.PrizeBunRepository{DB: db}
	inserted := 0
	for _, p := range prizes {
		if err := prizeRepo.Save(ctx, p); err != nil {
			slog.Error("failed to insert prize", "error", err)
			continue
		}
		inserted++
	}

	fmt.Printf("Inserted %d documents\n", inserted)
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session")
		if err != nil || cookie.Value == "" {
			// Si c'est une requête HTMX, on peut renvoyer un header pour rediriger côté client
			// ou simplement renvoyer vers la page de login
			if c.Request().Header.Get("HX-Request") == "true" {
				c.Response().Header().Set("HX-Redirect", "/login")
				return nil
			}
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		return next(c)
	}
}

func initWeb(userService *app.UserService, prizeService *app.PrizeService, authService *app.AuthService) *echo.Echo {
	handler := web.NewHandler(userService, prizeService, authService)

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				slog.Error("request error", "method", v.Method, "uri", v.URI, "status", v.Status, "error", v.Error)
			} else {
				slog.Info("request", "method", v.Method, "uri", v.URI, "status", v.Status)
			}
			return nil
		},
	}))
	e.Use(middleware.Gzip())

	e.GET(web.RouteIndex, handler.HandleIndexPage)
	e.GET(web.RoutePrize, handler.HandlePrizePage)
	e.GET(web.RouteAdmin, handler.HandleAdminPage, AuthMiddleware)
	e.GET(web.RouteAbout, handler.HandleAboutPage)
	e.GET(web.RouteLogin, handler.HandleLoginPage)
	e.POST(web.RouteLogin, handler.HandleLoginPost)
	e.POST(web.RouteLogout, handler.HandleLogout)
	e.POST(web.RouteSwitch, handler.HandleUserStatusSwitch, AuthMiddleware)
	e.GET(web.RouteStatus, func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Servir les fichiers statiques depuis le système de fichiers embarqué
	staticSubFS, err := fs.Sub(web.StaticFS, "static")
	if err != nil {
		slog.Error("Error creating sub-filesystem for static files", "error", err)
		os.Exit(1)
	}
	e.StaticFS(web.RouteStatic, staticSubFS)

	return e
}
