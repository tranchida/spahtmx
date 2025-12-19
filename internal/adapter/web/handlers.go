package web

import (
	"context"
	"io"
	"log"
	"net/http"
	"spahtmx/internal/adapter/web/templates"
	"spahtmx/internal/app"

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
	service app.UserService
}

func NewHandler(userService app.UserService) *Handler {
	return &Handler{
		service: userService,
	}
}

func (h *Handler) HandleIndexPage(c echo.Context) error {
	return handlePage(c, RouteIndex, templates.Index())
}

func (h *Handler) HandleAdminPage(c echo.Context) error {

	users := h.service.GetUsers(c.Request().Context())
	usersCount := h.service.GetUserCount(c.Request().Context())
	pageViews := h.service.GetPageView(c.Request().Context())

	return handlePage(c, RouteAdmin, templates.Admin(users, usersCount, pageViews))
}

func (h *Handler) HandleAboutPage(c echo.Context) error {
	return handlePage(c, RouteAbout, templates.About())
}

func (h *Handler) HandleUserStatusSwitch(c echo.Context) error {

	id := c.Param("id")
	h.service.UpdateUserStatus(c.Request().Context(), id)

	return handlePage(c, RouteAdmin, templates.Userlist(h.service.GetUsers(c.Request().Context())))
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
