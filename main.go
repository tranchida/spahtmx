package main

import (
	"log"
	"spahtmx/internal/adapter/gorm"
	"spahtmx/internal/adapter/web"
	"spahtmx/internal/app"
)



func main() {

	db := gorm.InitDB()

	repo := gorm.NewGormRepository(db)

	userService := app.NewUserService(repo)

	e := web.InitWeb(*userService)

	log.Println("🚀 Serveur démarré sur http://localhost:8765")
	e.Start(":8765")

}

