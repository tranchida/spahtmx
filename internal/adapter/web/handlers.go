package web

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"spahtmx/internal/adapter/web/templates"
	"spahtmx/internal/app"
	"spahtmx/internal/domain"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

const (
	RouteIndex  = "/"
	RouteAdmin  = "/admin"
	RouteAbout  = "/about"
	RouteStatus = "/status"
	RoutePrize  = "/prize"
	RouteSwitch = "/api/switch/:id"
	RouteStatic = "/static"
)

type Handler struct {
	userService  *app.UserService
	prizeService *app.PrizeService
}

func NewHandler(userService *app.UserService, prizeService *app.PrizeService) *Handler {
	return &Handler{
		userService:  userService,
		prizeService: prizeService,
	}
}

func (h *Handler) HandleIndexPage(c echo.Context) error {
	return handlePage(c, RouteIndex, templates.Index())
}

func (h *Handler) HandleAdminPage(c echo.Context) error {

	users, err := h.userService.GetUsers(c.Request().Context())
	if err != nil {
		return translateError(err)
	}
	usersCount := h.userService.GetUserCount(c.Request().Context())
	pageViews := h.userService.GetPageView(c.Request().Context())

	return handlePage(c, RouteAdmin, templates.Admin(users, usersCount, pageViews))
}

func (h *Handler) HandleAboutPage(c echo.Context) error {
	return handlePage(c, RouteAbout, templates.About())
}

func (h *Handler) HandleUserStatusSwitch(c echo.Context) error {

	id := c.Param("id")
	if err := h.userService.UpdateUserStatus(c.Request().Context(), id); err != nil {
		return translateError(err)
	}

	users, err := h.userService.GetUsers(c.Request().Context())
	if err != nil {
		return translateError(err)
	}

	return handlePage(c, RouteAdmin, templates.Userlist(users))
}

func (h *Handler) HandlePrizePage(c echo.Context) error {
	prizes, err := h.prizeService.GetPrizes(c.Request().Context())
	if err != nil {
		return translateError(err)
	}
	return handlePage(c, RoutePrize, templates.Prize(prizes))
}

func translateError(err error) error {
	if errors.Is(err, domain.ErrUserNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	if errors.Is(err, domain.ErrInvalidInput) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error").SetInternal(err)
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
		slog.Error("Render error", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Erreur de rendu").SetInternal(err)
	}
	return nil
}
