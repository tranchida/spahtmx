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

	model.ConnectDatabase()

	e := echo.New()
	e.GET(RouteIndex, handleIndexPage)
	e.GET(RouteAdmin, handleAdminPage)
	e.GET(RouteAbout, handleAboutPage)
	e.POST(RouteSwitch, handleUserStatusSwitch)
	e.GET(RouteStatus, func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Servir les fichiers statiques depuis le système de fichiers embarqué
	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("Erreur lors de la création du sous-système de fichiers: %v", err)
	}
	e.StaticFS(RouteStatic, staticSubFS)

	log.Println("🚀 Serveur démarré sur http://localhost:8765")
	e.Start(":8765")

}

func handleIndexPage(c echo.Context) error {
	return handlePage(c, RouteIndex, templates.Index())
}

func handleAdminPage(c echo.Context) error {

	users := model.GetUsers()
	usersCount := model.GetUserCount()
	pageViews := model.GetPageView()

	return handlePage(c, RouteAdmin, templates.Admin(users, usersCount, pageViews))
}

func handleAboutPage(c echo.Context) error {
	return handlePage(c, RouteAbout, templates.About())
}

func handleUserStatusSwitch(c echo.Context) error{

	id := c.Param("id")
	model.UpdateUserStatus(c.Request().Context(), id)

	return handlePage(c, RouteAdmin, templates.Userlist(model.GetUsers()))
}

func handlePage(c echo.Context, page string, contents templ.Component) error {
    fragment := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        if err := templates.Nav(page).Render(ctx, w); err != nil {
            return err
        }
        return contents.Render(ctx, w)
    })

    isHTMXRequest := c.Request().Header.Get("HX-Request") == "true"
    var component templ.Component
    if isHTMXRequest {
        component = fragment
    } else {
        component = templates.Base(page, contents)
    }

    if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
        log.Printf("Erreur lors du rendu: %v", err)
        http.Error(c.Response().Writer, "Erreur interne", http.StatusInternalServerError)
    }
    return nil
}
