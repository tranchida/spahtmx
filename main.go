package main

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"spahtmx/internal/model"
	"spahtmx/templates"
	"strconv"

	"github.com/a-h/templ"
)

//go:embed static/*
var staticFS embed.FS

func main() {

	model.ConnectDatabase()

	http.HandleFunc("/", handleIndexPage)

	http.HandleFunc("/admin", handleAdminPage)

	http.HandleFunc("/about", handleAboutPage)

	http.HandleFunc("/api/switch/{id}", handleUserStatusSwitch)

	// Servir les fichiers statiques depuis le système de fichiers embarqué
	staticSubFS, _ := fs.Sub(staticFS, "static")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSubFS))))

	log.Println("🚀 Serveur démarré sur http://localhost:8765")
	log.Fatal(http.ListenAndServe(":8765", nil))
}

func handleIndexPage(writer http.ResponseWriter, request *http.Request) {
	handlePage(writer, request, "/", templates.Index())
}

func handleAdminPage(writer http.ResponseWriter, request *http.Request) {

	users := model.GetUsers()
	usersCount := model.GetUserCount()
	pageViews := model.GetPageView()

	handlePage(writer, request, "/admin", templates.Admin(users, usersCount, pageViews))
}

func handleAboutPage(writer http.ResponseWriter, request *http.Request) {
	handlePage(writer, request, "/about", templates.About())
}

func handleUserStatusSwitch(writer http.ResponseWriter, request *http.Request) {

	id := request.PathValue("id")
	idval, err := strconv.Atoi(id)
	if err == nil {
		user := model.User{}
		model.DB.First(&user, idval)
		user.Status = !user.Status
		model.DB.Save(&user)
	}

	handlePage(writer, request, "/admin", templates.Userlist(model.GetUsers()))
}

func handlePage(writer http.ResponseWriter, request *http.Request, page string, contents templ.Component) {
	// Détecter si c'est une requête HTMX via le header

	if request.Header.Get("HX-Request") == "true" {
		fragment := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if err := templates.Nav(page).Render(ctx, w); err != nil {
				return err
			}
			return contents.Render(ctx, w)
		})
		err := fragment.Render(request.Context(), writer)
		if err != nil {
			log.Printf("Erreur lors du rendu du fragment: %v", err)
			http.Error(writer, "Erreur interne du serveur", http.StatusInternalServerError)
		}
	} else {
		err := templates.Base(page, contents).Render(request.Context(), writer)
		if err != nil {
			log.Printf("Erreur lors du rendu du template: %v", err)
			http.Error(writer, "Erreur interne du serveur", http.StatusInternalServerError)
		}
	}
}
