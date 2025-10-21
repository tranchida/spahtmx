# SPA HTMX avec Go

Une application Single Page Application (SPA) moderne utilisant HTMX, Tailwind CSS et les templates Go.

## 🚀 Fonctionnalités

- **3 pages** : Accueil, Admin, À propos
- **Navigation SPA** : Pas de rechargement de page grâce à HTMX
- **Templates Go** : Rendu côté serveur avec le système de templates natif de Go
- **Design moderne** : Interface responsive avec Tailwind CSS et animations fluides
- **Icônes SVG** : Interface enrichie avec des icônes intégrées

## 📁 Structure du projet

```
.
├── main.go              # Serveur HTTP et handlers
├── templates/           # Templates HTML
│   ├── base.html       # Template de base avec navigation
│   ├── index.html      # Page d'accueil
│   ├── admin.html      # Page admin avec statistiques
│   └── about.html      # Page à propos
└── static/             # Fichiers statiques
    ├── css/
    │   └── style.css   # (non utilisé - remplacé par Tailwind CSS)
    └── js/             # Scripts JavaScript
```

## 🛠️ Installation et démarrage

1. Assurez-vous d'avoir Go installé (version 1.25+)

2. Clonez le projet et accédez au répertoire :
```bash
cd /home/xcigta/dev/test/web/htmx/spahtmx
```

3. Lancez le serveur :
```bash
go run main.go
```

4. Ouvrez votre navigateur à l'adresse : **http://localhost:8765**

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
- `/admin` - Page admin complète avec statistiques
- `/about` - Page à propos complète

**Fragments HTMX** (navigation SPA) :
- `/page/index` - Fragment de contenu pour l'accueil
- `/page/admin` - Fragment de contenu pour admin
- `/page/about` - Fragment de contenu pour à propos

**Fichiers statiques** :
- `/static/*` - Serveur de fichiers statiques

## 🎨 Personnalisation

### Modifier les templates
Les templates HTML se trouvent dans le dossier `templates/` :
- `base.html` : Layout principal avec navigation et configuration Tailwind
- `index.html`, `admin.html`, `about.html` : Contenu des pages

### Personnaliser les couleurs Tailwind
Dans `templates/base.html`, modifiez la configuration Tailwind :
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
1. Créez un nouveau template dans `templates/` (ex: `contact.html`)
2. Ajoutez les routes dans `main.go` :
```go
http.HandleFunc("/contact", handleContact)
http.HandleFunc("/page/contact", handleContactFragment)
```
3. Implémentez les handlers correspondants

## 📝 Technologies

- **Go 1.25** - Backend et serveur HTTP natif
- **HTMX 1.9.10** - Interactions AJAX sans JavaScript complexe
- **Tailwind CSS** - Framework CSS utility-first via CDN
- **HTML Templates** - Système de templates natif de Go
- **SVG Icons** - Icônes vectorielles intégrées

## 🌟 Avantages de cette stack

- ✅ **Simplicité** : Pas de build frontend, pas de npm, pas de Node.js
- ✅ **Performance** : Serveur Go ultra-rapide et léger
- ✅ **SEO-friendly** : Rendu côté serveur pour toutes les pages
- ✅ **Expérience utilisateur** : Navigation fluide comme une SPA React
- ✅ **Maintenabilité** : Code simple et facile à comprendre
- ✅ **Production-ready** : Binaire Go compilé, facile à déployer

## 🚀 Déploiement

### Compilation
```bash
go build -o spahtmx main.go
```

### Exécution en production
```bash
./spahtmx
```

Le serveur écoute sur le port 8765. Vous pouvez modifier ce port dans `main.go` si nécessaire.

## 📄 License

Projet d'exemple - Libre d'utilisation
