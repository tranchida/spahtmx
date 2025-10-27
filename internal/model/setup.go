package model

import (
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

	var count int64
	DB.Model(&User{}).Count(&count)
	log.Printf("Count %d", count)
	if count == 0 {
		users = []User{
			{ID: 1, Username: "alice", Email: "alice@fake.com", Status: true},
			{ID: 2, Username: "bob", Email: "bob@fake.com", Status: false},
			{ID: 3, Username: "charlie", Email: "charlie@fake.com", Status: true},
		}
		DB.Create(users).Save(users)
	}

}
