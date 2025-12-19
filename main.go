package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"spahtmx/internal/adapter/mongodb"
	"spahtmx/internal/adapter/web"
	"spahtmx/internal/app"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

//go:embed static/*
var staticFS embed.FS

func main() {

	db := initDB()

	repo := mongodb.MongoRepository{
		DB: db,
	}

	userService := app.NewUserService(repo)

	e := initWeb(*userService)

	log.Println("🚀 Serveur démarré sur http://localhost:8765")
	err := e.Start(":8765")
	if err != nil {
		log.Fatal("Server start failed", err)
	}

}

func initDB() *mongo.Database {

	client, err := mongo.Connect(
		options.Client().ApplyURI("mongodb://localhost:27017"),
		options.Client().SetAuth(options.Credential{
			Username: "root",
			Password: "example",
		}),
	)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}

	us := []mongodb.UserMongo{
		{ID: bson.NewObjectID(), Username: "alice", Email: "alice@fake.com", Status: true},
		{ID: bson.NewObjectID(), Username: "bob", Email: "bob@fake.com", Status: false},
		{ID: bson.NewObjectID(), Username: "charlie", Email: "charlie@fake.com", Status: true},
	}

	err = client.Database("test").Drop(context.Background())
	if err != nil {
		return nil
	}

	_, err = client.Database("test").Collection("users").InsertMany(context.Background(), us)
	if err != nil {
		log.Fatal("Failed to insert users", err)
	}

	return client.Database("test")
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
