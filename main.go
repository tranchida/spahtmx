package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"spahtmx/internal/adapter/gorm"
	"spahtmx/internal/adapter/web"
	"spahtmx/internal/app"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"

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
