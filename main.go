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
	Page      string
	Users     []model.User
	UserCount int
	PageViews int
}

var templates *template.Template

func init() {
	// Charger tous les templates
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	// Routes pour les pages complètes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/admin", handleAdmin)
	http.HandleFunc("/about", handleAbout)

	// Routes pour les fragments HTMX
	http.HandleFunc("/page/index", handleIndexFragment)
	http.HandleFunc("/page/admin", handleAdminFragment)
	http.HandleFunc("/page/about", handleAboutFragment)

	http.HandleFunc("/api/switch/{id}", handleUserStatusSwitch)
	// Servir les fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("🚀 Serveur démarré sur http://localhost:8765")
	log.Fatal(http.ListenAndServe(":8765", nil))
}

func handleUserStatusSwitch(writer http.ResponseWriter, request *http.Request) {

	id := request.PathValue("id")
	idval, err := strconv.Atoi(id)
	if err == nil {
		model.GetUsers()[idval-1].Status = !model.GetUsers()[idval-1].Status
	}

	data := PageData{
		Page:  "userlist",
		Users: model.GetUsers(),
	}

	err = templates.ExecuteTemplate(writer, "userlist", data)

	if err != nil {
		log.Printf("Erreur lors du rendu du fragment: %v", err)
		http.Error(writer, "Erreur interne du serveur", http.StatusInternalServerError)
	}

}

// Handlers pour les pages complètes (première visite)
func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := PageData{Title: "Accueil - SPA HTMX"}
	renderFullPage(w, data)
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:     "Admin - SPA HTMX",
		Page:      "admin",
		UserCount: 142,
		PageViews: 3789,
		Users:     model.GetUsers(),
	}
	renderFullPage(w, data)
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "À propos - SPA HTMX",
		Page:  "about",
	}
	renderFullPage(w, data)
}

// Handlers pour les fragments HTMX (navigation SPA)
func handleIndexFragment(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Accueil - SPA HTMX",
		Page:  "index",
	}
	renderFragment(w, data)
}

func handleAdminFragment(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:     "Admin - SPA HTMX",
		Page:      "admin",
		UserCount: 142,
		PageViews: 3789,
		Users:     model.GetUsers(),
	}
	renderFragment(w, data)
}

func handleAboutFragment(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "À propos - SPA HTMX",
		Page:  "about",
	}
	renderFragment(w, data)
}

// Fonction pour rendre une page complète
func renderFullPage(w http.ResponseWriter, data PageData) {

	err := templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}

// Fonction pour rendre uniquement le contenu (pour HTMX)
func renderFragment(w http.ResponseWriter, data PageData) {

	err := templates.ExecuteTemplate(w, data.Page, data)
	if err != nil {
		log.Printf("Erreur lors du rendu du fragment: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}
