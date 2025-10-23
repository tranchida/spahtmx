package main

import (
	"embed"
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
	handlePage(writer, request, templates.Index())
}

func handleAdminPage(writer http.ResponseWriter, request *http.Request) {

	users := model.GetUsers()
	usersCount := model.GetUserCount()
	pageViews := model.GetPageView()

	handlePage(writer, request, templates.Admin(users, usersCount, pageViews))
}

func handleAboutPage(writer http.ResponseWriter, request *http.Request) {
	handlePage(writer, request, templates.About())
}

func handleUserStatusSwitch(writer http.ResponseWriter, request *http.Request) {

	id := request.PathValue("id")
	idval, err := strconv.Atoi(id)
	if err == nil {
		model.GetUsers()[idval-1].Status = !model.GetUsers()[idval-1].Status
	}

	handlePage(writer, request, templates.Userlist(model.GetUsers()))
}

func handlePage(writer http.ResponseWriter, request *http.Request, contents templ.Component) {
	// Détecter si c'est une requête HTMX via le header
	if request.Header.Get("HX-Request") == "true" {
		err := contents.Render(request.Context(), writer)
		if err != nil {
			log.Printf("Erreur lors du rendu du fragment: %v", err)
			http.Error(writer, "Erreur interne du serveur", http.StatusInternalServerError)
		}
	} else {
		err := templates.Base(contents).Render(request.Context(), writer)
		if err != nil {
			log.Printf("Erreur lors du rendu du template: %v", err)
			http.Error(writer, "Erreur interne du serveur", http.StatusInternalServerError)
		}
	}
}
