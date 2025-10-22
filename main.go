package main

import (
	"html/template"
	"log"
	"net/http"
	"spahtmx/internal/model"
	"strconv"
)

type PageData struct {
	Title     string
	Users     []model.User
	UserCount int
	PageViews int
}

var templates *template.Template

func init() {
	// Charger tous les templates
	templates = template.Must(template.ParseFiles("templates/base.html", "templates/userlist.html"))
}

func main() {

	http.HandleFunc("/", handlePage("index", func() PageData {
		return PageData{
			Title: "Accueil - SPA HTMX",
		}
	}))

	http.HandleFunc("/admin", handlePage("admin", func() PageData {
		return PageData{
			Title:     "Admin - SPA HTMX",
			UserCount: 142,
			PageViews: 3789,
			Users:     model.GetUsers(),
		}
	}))

	http.HandleFunc("/about", handlePage("about", func() PageData {
		return PageData{
			Title: "À propos - SPA HTMX",
		}
	}))

	http.HandleFunc("/api/switch/{id}", handleUserStatusSwitch)

	// Servir les fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("🚀 Serveur démarré sur http://localhost:8765")
	log.Fatal(http.ListenAndServe(":8765", nil))
}

func handlePage(page string, dataFunc func() PageData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := dataFunc()
		// Détecter si c'est une requête HTMX via le header
		if r.Header.Get("HX-Request") == "true" {
			renderFragment(w, page, data)
		} else {
			renderFullPage(w, page, data)
		}
	}
}

func handleUserStatusSwitch(writer http.ResponseWriter, request *http.Request) {

	id := request.PathValue("id")
	idval, err := strconv.Atoi(id)
	if err == nil {
		model.GetUsers()[idval-1].Status = !model.GetUsers()[idval-1].Status
	}

	data := PageData{
		Users: model.GetUsers(),
	}

	tmpl := template.Must(templates.Clone())
	err = tmpl.ExecuteTemplate(writer, "userlist", data)

	if err != nil {
		log.Printf("Erreur lors du rendu du fragment: %v", err)
		http.Error(writer, "Erreur interne du serveur", http.StatusInternalServerError)
	}

}

// Fonction pour rendre une page complète
func renderFullPage(w http.ResponseWriter, page string, data PageData) {

	tmpl := template.Must(template.Must(templates.Clone()).ParseFiles("templates/" + page + ".html"))
	err := tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}

// Fonction pour rendre uniquement le contenu (pour HTMX)
func renderFragment(w http.ResponseWriter, page string, data PageData) {

	tmpl := template.Must(template.Must(templates.Clone()).ParseFiles("templates/" + page + ".html"))
	err := tmpl.ExecuteTemplate(w, "content", data)
	if err != nil {
		log.Printf("Erreur lors du rendu du fragment: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}
