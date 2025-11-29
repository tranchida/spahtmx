# SPA HTMX avec Go

Une application Single Page Application (SPA) moderne utilisant HTMX, Tailwind CSS et Templ.

## 🚀 Fonctionnalités

- **3 pages** : Accueil, Admin, À propos
- **Navigation SPA** : Pas de rechargement de page grâce à HTMX
- **Templates Templ** : Rendu côté serveur avec Templ (type-safe Go templates)
- **Design moderne** : Interface responsive avec Tailwind CSS et animations fluides
- **Gestion d'utilisateurs** : Page admin avec liste d'utilisateurs et statistiques
- **API interactive** : Toggle du statut utilisateur avec HTMX
- **Fichiers statiques embarqués** : Déploiement simplifié avec embed.FS

## 📁 Structure du projet

```
.
├── main.go              # Serveur HTTP et handlers
├── go.mod               # Dépendances Go (templ, air)
├── internal/
│   └── model/
│       └── user.go      # Modèle User et fonctions de données
├── templates/           # Templates Templ
│   ├── base.templ       # Template de base avec navigation
│   ├── nav.templ        # Composant navigation
│   ├── footer.templ     # Composant footer
│   ├── index.templ      # Page d'accueil
│   ├── admin.templ      # Page admin avec statistiques
│   ├── about.templ      # Page à propos
│   ├── userlist.templ   # Liste d'utilisateurs
│   └── *_templ.go       # Fichiers générés par Templ
└── static/              # Fichiers statiques (embarqués)
    └── js/
        ├── htmx.min.js
        └── tailwind.min.js
```

## 🛠️ Installation et démarrage

1. Assurez-vous d'avoir Go installé (version 1.25+)

2. Clonez le projet et accédez au répertoire :
```bash
cd spahtmx
```

3. Installez les dépendances :
```bash
go mod download
```

4. Générez les templates Templ (si modifiés) :
```bash
go tool templ generate
```

5. Lancez le serveur :
```bash
go run main.go
```
Ou utilisez Air pour le développement avec rechargement automatique :
```bash
go tool air
```

6. Ouvrez votre navigateur à l'adresse : **http://localhost:8765**

## 🎯 Comment ça fonctionne

### Architecture SPA avec HTMX

L'application utilise HTMX pour créer une expérience SPA sans framework JavaScript lourd :

- **Première visite** : Le serveur envoie la page HTML complète avec le layout (navbar, footer, contenu)
- **Navigation** : Les clics sur les liens déclenchent des requêtes HTMX qui ne chargent que le contenu de la page
- **Historique** : HTMX gère automatiquement l'historique du navigateur avec `hx-push-url`
- **Transitions** : Animations CSS fluides lors des changements de page

### Routes

**Pages complètes** (première visite ou rechargement) :
- `/` - Page d'accueil complète
- `/admin` - Page admin complète avec statistiques et liste d'utilisateurs
- `/about` - Page à propos complète

**API HTMX** :
- `/api/switch/{id}` - Toggle du statut d'un utilisateur (retourne la liste mise à jour)

**Fichiers statiques** :
- `/static/*` - Serveur de fichiers statiques (embarqués avec embed.FS)

## 🎨 Personnalisation

Lance le serveur avec Air et recompilation automatique des templates
```Bash
make dev 
```

### Modifier les templates
Les templates Templ se trouvent dans le dossier `templates/` avec l'extension `.templ` :
- `base.templ` : Layout principal avec navigation et configuration Tailwind
- `nav.templ`, `footer.templ` : Composants de navigation et footer
- `index.templ`, `admin.templ`, `about.templ` : Contenu des pages
- `userlist.templ` : Composant de liste d'utilisateurs

### Personnaliser les couleurs Tailwind
Dans `templates/base.templ`, modifiez la configuration Tailwind :
```javascript
tailwind.config = {
    theme: {
        extend: {
            colors: {
                primary: '#667eea',    // Couleur principale
                secondary: '#764ba2',  // Couleur secondaire
            }
        }
    }
}
```

### Ajouter de nouvelles pages
1. Créez un nouveau template dans `templates/` (ex: `contact.templ`)
```go
package templates

templ Contact() {
    <div class="bg-white rounded-xl shadow-lg p-8">
        <h1 class="text-3xl font-bold mb-4">Contact</h1>
        // Votre contenu ici
    </div>
}
```
2. Ajoutez la route dans `main.go` :
```go
http.HandleFunc("/contact", handleContactPage)
```
3. Implémentez le handler :
```go
func handleContactPage(writer http.ResponseWriter, request *http.Request) {
    handlePage(writer, request, templates.Contact())
}
```

## 📝 Technologies

- **Go 1.25** - Backend et serveur HTTP natif
- **Templ** - Templates type-safe pour Go (github.com/a-h/templ)
- **HTMX** - Interactions AJAX sans JavaScript complexe
- **Tailwind CSS** - Framework CSS utility-first via CDN
- **Air** - Rechargement automatique pour le développement
- **embed.FS** - Fichiers statiques embarqués dans le binaire

## 🌟 Avantages de cette stack

- ✅ **Simplicité** : Pas de build frontend complexe, pas de npm massif
- ✅ **Type-safety** : Templ fournit des templates type-safe avec autocomplétion
- ✅ **Performance** : Serveur Go ultra-rapide et léger
- ✅ **SEO-friendly** : Rendu côté serveur pour toutes les pages
- ✅ **Expérience utilisateur** : Navigation fluide comme une SPA React
- ✅ **Maintenabilité** : Code Go pur, facile à comprendre et déboguer
- ✅ **Production-ready** : Binaire unique avec assets embarqués, déploiement simple
- ✅ **Hot reload** : Développement rapide avec Air

## 🚀 Déploiement

### Compilation
```bash
go build -o spahtmx main.go
```

### Exécution en production
```bash
./spahtmx
```

Le serveur écoute sur le port 8080. Vous pouvez modifier ce port dans `main.go` si nécessaire.

## 📄 License

Projet d'exemple - Libre d'utilisation
