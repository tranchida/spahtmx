package web

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sort"
	"spahtmx/internal/adapter/web/templates"
	"spahtmx/internal/app"
	"spahtmx/internal/domain"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

const (
	RouteIndex  = "/"
	RouteAdmin  = "/admin"
	RouteAbout  = "/about"
	RouteStatus = "/status"
	RoutePrize  = "/prize"
	RouteLogin  = "/login"
	RouteLogout = "/logout"
	RouteSwitch = "/api/switch/:id"
	RouteStatic = "/static"
)

type Handler struct {
	userService  *app.UserService
	prizeService *app.PrizeService
	authService  *app.AuthService
}

func NewHandler(userService *app.UserService, prizeService *app.PrizeService, authService *app.AuthService) *Handler {
	return &Handler{
		userService:  userService,
		prizeService: prizeService,
		authService:  authService,
	}
}

func (h *Handler) HandleLoginPage(c echo.Context) error {
	return h.handlePage(c, RouteLogin, templates.Login(""))
}

func (h *Handler) HandleLoginPost(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.authService.Login(c.Request().Context(), username, password)
	if err != nil {
		return h.handlePage(c, RouteLogin, templates.Login("Identifiants incorrects"))
	}

	// Création du cookie de session (très basique pour l'exemple)
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = user.Username
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	// On stocke l'utilisateur dans le contexte pour handlePage
	c.Set("user", user)

	return h.HandleIndexPage(c)
}

func (h *Handler) HandleLogout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	c.SetCookie(cookie)

	// On marque explicitement qu'il n'y a plus d'utilisateur pour handlePage
	c.Set("user", nil)
	c.Set("logout", true)

	return h.HandleIndexPage(c)
}

func (h *Handler) HandleIndexPage(c echo.Context) error {
	return h.handlePage(c, RouteIndex, templates.Index())
}

func (h *Handler) HandleAdminPage(c echo.Context) error {

	users, err := h.userService.GetUsers(c.Request().Context())
	if err != nil {
		return translateError(err)
	}
	usersCount := h.userService.GetUserCount(c.Request().Context())
	pageViews := h.userService.GetPageView(c.Request().Context())

	return h.handlePage(c, RouteAdmin, templates.Admin(users, usersCount, pageViews))
}

func (h *Handler) HandleAboutPage(c echo.Context) error {
	return h.handlePage(c, RouteAbout, templates.About())
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

	return h.handlePage(c, RouteAdmin, templates.Userlist(users))
}

func (h *Handler) HandlePrizePage(c echo.Context) error {
	category := c.QueryParam("category")
	year := c.QueryParam("year")

	var prizes []domain.Prize
	var err error

	if category != "" && year != "" {
		prizes, err = h.prizeService.GetPrizesByCategoryAndYear(c.Request().Context(), category, year)
	} else if category != "" {
		prizes, err = h.prizeService.GetPrizesByCategory(c.Request().Context(), category)
	} else if year != "" {
		prizes, err = h.prizeService.GetPrizesByYear(c.Request().Context(), year)
	} else {
		currentYear := strconv.Itoa(time.Now().Year())
		prizes, err = h.prizeService.GetPrizesByYear(c.Request().Context(), currentYear)
	}

	if err != nil {
		return translateError(err)
	}

	categories, err := h.prizeService.GetCategories(c.Request().Context())
	if err != nil {
		return translateError(err)
	}
	years, err := h.prizeService.GetYears(c.Request().Context())
	if err != nil {
		return translateError(err)
	}

	sort.Strings(categories)
	sort.Slice(years, func(i, j int) bool { return years[i] > years[j] })

	return h.handlePage(c, RoutePrize, templates.Prize(prizes, categories, years, category, year))
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

func (h *Handler) handlePage(c echo.Context, page string, contents templ.Component) error {
	var user *domain.User

	// On vérifie d'abord si l'utilisateur est dans le contexte (cas du login/logout)
	if u, ok := c.Get("user").(domain.User); ok {
		user = &u
	} else if uPtr, ok := c.Get("user").(*domain.User); ok {
		user = uPtr
	} else if c.Get("logout") == nil {
		// Sinon on cherche dans le cookie, sauf si on vient de se déconnecter
		if cookie, err := c.Cookie("session"); err == nil && cookie.Value != "" {
			if u, err := h.authService.GetUserByUsername(c.Request().Context(), cookie.Value); err == nil {
				user = &u
			}
		}
	}

	fragment := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := templates.Nav(page, user).Render(ctx, w); err != nil {
			return err
		}
		return contents.Render(ctx, w)
	})

	isHTMXRequest := c.Request().Header.Get("HX-Request") == "true"
	var component templ.Component
	if isHTMXRequest {
		component = fragment
	} else {
		component = templates.Base(page, user, contents)
	}

	if err := component.Render(c.Request().Context(), c.Response().Writer); err != nil {
		slog.Error("Render error", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Erreur de rendu").SetInternal(err)
	}
	return nil
}
