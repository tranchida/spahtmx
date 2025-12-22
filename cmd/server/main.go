package main

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"spahtmx/internal/adapter/mongodb"
	"spahtmx/internal/adapter/web"
	"spahtmx/internal/app"
	"spahtmx/internal/config"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, client := initDB(ctx, cfg)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			slog.Error("Error disconnecting from MongoDB", "error", err)
		}
	}()

	repo := mongodb.MongoRepository{
		DB: db,
	}

	userService := app.NewUserService(repo)

	e := initWeb(*userService)

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

func initDB(ctx context.Context, cfg *config.Config) (*mongo.Database, *mongo.Client) {
	client, err := mongo.Connect(
		options.Client().ApplyURI(cfg.MongoDBURL),
	)
	if err != nil {
		slog.Error("Failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}

	// Ping pour vérifier la connexion
	if err := client.Ping(ctx, nil); err != nil {
		slog.Error("Failed to ping MongoDB", "error", err)
		os.Exit(1)
	}

	db := client.Database("test")

	// Seed data (Optionnel pour le dev)
	if cfg.SeedDB {
		seedDatabase(ctx, db)
	}

	return db, client
}

func seedDatabase(ctx context.Context, db *mongo.Database) {
	us := []mongodb.UserMongo{
		{ID: bson.NewObjectID(), Username: "alice", Email: "alice@fake.com", Status: true},
		{ID: bson.NewObjectID(), Username: "bob", Email: "bob@fake.com", Status: false},
		{ID: bson.NewObjectID(), Username: "charlie", Email: "charlie@fake.com", Status: true},
	}

	if err := db.Drop(ctx); err != nil {
		slog.Warn("Failed to drop database", "error", err)
	}

	_, err := db.Collection("users").InsertMany(ctx, us)
	if err != nil {
		slog.Error("Failed to insert users", "error", err)
		os.Exit(1)
	}
	slog.Info("Database seeded successfully")
}

func initWeb(userService app.UserService) *echo.Echo {
	handler := web.NewHandler(userService)

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
	e.GET(web.RouteAdmin, handler.HandleAdminPage)
	e.GET(web.RouteAbout, handler.HandleAboutPage)
	e.POST(web.RouteSwitch, handler.HandleUserStatusSwitch)
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
