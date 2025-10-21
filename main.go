package main

import (
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Title     string
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

	// Servir les fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("🚀 Serveur démarré sur http://localhost:8765")
	log.Fatal(http.ListenAndServe(":8765", nil))
}

// Handlers pour les pages complètes (première visite)
func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := PageData{Title: "Accueil - SPA HTMX"}
	renderFullPage(w, "index.html", data)
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:     "Admin - SPA HTMX",
		UserCount: 142,
		PageViews: 3789,
	}
	renderFullPage(w, "admin.html", data)
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "À propos - SPA HTMX"}
	renderFullPage(w, "about.html", data)
}

// Handlers pour les fragments HTMX (navigation SPA)
func handleIndexFragment(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Accueil - SPA HTMX"}
	renderFragment(w, "index.html", data)
}

func handleAdminFragment(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:     "Admin - SPA HTMX",
		UserCount: 142,
		PageViews: 3789,
	}
	renderFragment(w, "admin.html", data)
}

func handleAboutFragment(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "À propos - SPA HTMX"}
	renderFragment(w, "about.html", data)
}

// Fonction pour rendre une page complète
func renderFullPage(w http.ResponseWriter, templateName string, data PageData) {
	// Parse le template de base et le template de contenu ensemble
	tmpl, err := template.ParseFiles("templates/base.html", "templates/"+templateName)
	if err != nil {
		log.Printf("Erreur lors du parsing du template: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}

// Fonction pour rendre uniquement le contenu (pour HTMX)
func renderFragment(w http.ResponseWriter, templateName string, data PageData) {
	tmpl, err := template.ParseFiles("templates/" + templateName)
	if err != nil {
		log.Printf("Erreur lors du parsing du fragment: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "content", data)
	if err != nil {
		log.Printf("Erreur lors du rendu du fragment: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}
