package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"spahtmx/internal/adapter/gorm"
	"spahtmx/internal/adapter/web"
	"spahtmx/internal/app"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
    RouteIndex  = "/"
    RouteAdmin  = "/admin"
    RouteAbout  = "/about"
    RouteStatus = "/status"
    RouteSwitch = "/api/switch/:id"
    RouteStatic = "/static"
)

//go:embed static/*
var staticFS embed.FS

func main() {

	db := gorm.InitDB()

	repo := gorm.NewGormRepository(db)

	userService := app.NewUserService(repo)

	handler := web.NewHandler(userService)

	e := echo.New()
	e.Use(middleware.Logger())
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

	log.Println("🚀 Serveur démarré sur http://localhost:8765")
	e.Start(":8765")

}

