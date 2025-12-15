package web

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"spahtmx/internal/app"
	"spahtmx/internal/adapter/web/templates"

	"github.com/a-h/templ"
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

type Handler struct {
	model app.UserService
}

func InitWeb(userService app.UserService) *echo.Echo{

    handler := NewHandler(userService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.GET(RouteIndex, handler.HandleIndexPage)
	e.GET(RouteAdmin, handler.HandleAdminPage)
	e.GET(RouteAbout, handler.HandleAboutPage)
	e.POST(RouteSwitch, handler.HandleUserStatusSwitch)
	e.GET(RouteStatus, func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Servir les fichiers statiques depuis le système de fichiers embarqué
	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("Erreur lors de la création du sous-système de fichiers: %v", err)
	}
	e.StaticFS(RouteStatic, staticSubFS)

    return e
}

func NewHandler(model app.UserService) *Handler {
	return &Handler{
		model: model,
	}
}

func (h *Handler) HandleIndexPage(c echo.Context) error {
	return handlePage(c, RouteIndex, templates.Index())
}

func (h *Handler) HandleAdminPage(c echo.Context) error {

	users := h.model.GetUsers()
	usersCount := h.model.GetUserCount()
	pageViews := h.model.GetPageView()

	return handlePage(c, RouteAdmin, templates.Admin(users, usersCount, pageViews))
}

func (h *Handler) HandleAboutPage(c echo.Context) error {
	return handlePage(c, RouteAbout, templates.About())
}

func (h *Handler) HandleUserStatusSwitch(c echo.Context) error{

	id := c.Param("id")
	h.model.UpdateUserStatus(c.Request().Context(), id)

	return handlePage(c, RouteAdmin, templates.Userlist(h.model.GetUsers()))
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
