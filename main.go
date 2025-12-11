package main

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"github.com/labstack/echo/v4"
	"spahtmx/internal/model"
	"spahtmx/templates"

	"github.com/a-h/templ"
)

//go:embed static/*
var staticFS embed.FS

func main() {

	model.ConnectDatabase()

	e := echo.New()
	e.GET("/", handleIndexPage)
	e.GET("/admin", handleAdminPage)
	e.GET("/about", handleAboutPage)
	e.POST("/api/switch/:id", handleUserStatusSwitch)
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Servir les fichiers statiques depuis le système de fichiers embarqué
	staticSubFS, _ := fs.Sub(staticFS, "static")
	e.StaticFS("/static", staticSubFS)

	log.Println("🚀 Serveur démarré sur http://localhost:8765")
	e.Start(":8765")

}

func handleIndexPage(c echo.Context) error {
	return handlePage(c, "/", templates.Index())
}

func handleAdminPage(c echo.Context) error {

	users := model.GetUsers()
	usersCount := model.GetUserCount()
	pageViews := model.GetPageView()

	return handlePage(c, "/admin", templates.Admin(users, usersCount, pageViews))
}

func handleAboutPage(c echo.Context) error {
	return handlePage(c, "/about", templates.About())
}

func handleUserStatusSwitch(c echo.Context) error{

	id := c.Param("id")
	model.UpdateUserStatus(c.Request().Context(), id)

	return handlePage(c, "/admin", templates.Userlist(model.GetUsers()))
}

func handlePage(c echo.Context, page string, contents templ.Component) error{
	// Détecter si c'est une requête HTMX via le header

	if c.Request().Header.Get("HX-Request") == "true" {
		fragment := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if err := templates.Nav(page).Render(ctx, w); err != nil {
				return err
			}
			return contents.Render(ctx, w)
		})
		err := fragment.Render(c.Request().Context(), c.Response().Writer)
		if err != nil {
			log.Printf("Erreur lors du rendu du fragment: %v", err)
			http.Error(c.Response().Writer, "Erreur interne du serveur", http.StatusInternalServerError)
		}
	} else {
		err := templates.Base(page, contents).Render(c.Request().Context(), c.Response().Writer)
		if err != nil {
			log.Printf("Erreur lors du rendu du template: %v", err)
			http.Error(c.Response().Writer, "Erreur interne du serveur", http.StatusInternalServerError)
		}
	}
	return nil
}
