package web

import (
	"context"
	"io"
	"log"
	"net/http"
	"spahtmx/internal/domain"
	"spahtmx/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

const (
    RouteIndex  = "/"
    RouteAdmin  = "/admin"
    RouteAbout  = "/about"
    RouteStatus = "/status"
    RouteSwitch = "/api/switch/:id"
    RouteStatic = "/static"
)

type Handler struct {
	model domain.UserRepository
}

func NewHandler(model domain.UserRepository) *Handler {
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
