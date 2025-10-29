package model

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&User{})
	if err != nil {
		return
	}

	DB = database

	count, err := gorm.G[User](DB).Count(context.Background(), "id")
	if err != nil {
		log.Fatal("count failed", err)
	}
	log.Printf("Count %d", count)
	if count == 0 {
		users := []User{
			{Username: "alice", Email: "alice@fake.com", Status: true},
			{Username: "bob", Email: "bob@fake.com", Status: false},
			{Username: "charlie", Email: "charlie@fake.com", Status: true},
		}
		err := gorm.G[[]User](DB).Create(context.Background(), &users)
		if err != nil {
			log.Fatal("insert seed records failed", err)
		}

	}

}
