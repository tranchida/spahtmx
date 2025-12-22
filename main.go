package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"spahtmx/internal/adapter/mongodb"
	"spahtmx/internal/adapter/web"
	"spahtmx/internal/app"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

//go:embed static/*
var staticFS embed.FS

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, client := initDB(ctx)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	repo := mongodb.MongoRepository{
		DB: db,
	}

	userService := app.NewUserService(repo)

	e := initWeb(*userService)

	// Démarrage du serveur dans une goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8765"
		}
		log.Printf("🚀 Serveur démarré sur http://localhost:%s", port)
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server start failed: %v", err)
		}
	}()

	// Attente du signal d'arrêt
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Arrêt gracieux du serveur Web avec un timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func initDB(ctx context.Context) (*mongo.Database, *mongo.Client) {
	mongoDBUrl := os.Getenv("MONGODB_URL")
	if mongoDBUrl == "" {
		mongoDBUrl = "mongodb://root:example@localhost:27017"
	}

	client, err := mongo.Connect(
		options.Client().ApplyURI(mongoDBUrl),
	)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	// Ping pour vérifier la connexion
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}

	db := client.Database("test")

	// Seed data (Optionnel pour le dev)
	if os.Getenv("SEED_DB") == "true" {
		us := []mongodb.UserMongo{
			{ID: bson.NewObjectID(), Username: "alice", Email: "alice@fake.com", Status: true},
			{ID: bson.NewObjectID(), Username: "bob", Email: "bob@fake.com", Status: false},
			{ID: bson.NewObjectID(), Username: "charlie", Email: "charlie@fake.com", Status: true},
		}

		if err := db.Drop(ctx); err != nil {
			log.Printf("Warning: failed to drop database: %v", err)
		}

		_, err = db.Collection("users").InsertMany(ctx, us)
		if err != nil {
			log.Fatal("Failed to insert users: ", err)
		}
	}

	return db, client
}

func initWeb(userService app.UserService) *echo.Echo {

	handler := web.NewHandler(userService)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Gzip())
	e.GET(web.RouteIndex, handler.HandleIndexPage)
	e.GET(web.RouteAdmin, handler.HandleAdminPage)
	e.GET(web.RouteAbout, handler.HandleAboutPage)
	e.POST(web.RouteSwitch, handler.HandleUserStatusSwitch)
	e.GET(web.RouteStatus, func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Servir les fichiers statiques depuis le système de fichiers embarqué
	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("Erreur lors de la création du sous-système de fichiers: %v", err)
	}
	e.StaticFS(web.RouteStatic, staticSubFS)

	return e
}
