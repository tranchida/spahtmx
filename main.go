package main

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"spahtmx/internal/model"
	"spahtmx/templates"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Serveur démarré sur http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
	model.UpdateUserStatus(request.Context(), id)

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
