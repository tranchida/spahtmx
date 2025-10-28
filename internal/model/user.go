package model

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	ID       int
	Username string
	Email    string
	Status   bool
}

func GetUsers() []User {

	var users []User
	DB.Find(&users)
	return users
}

func UpdateUserStatus(id string) {

	genericDB := gorm.G[User](DB)
	user, err := genericDB.Where("id = ?", id).Take(context.Background())
	if err != nil {
		return
	}
	_, err = genericDB.Where("id = ?", id).Update(context.Background(), "status", !user.Status)
	if err != nil {
		return
	}
}

func GetUserCount() string {
	return "210"
}

func GetPageView() string {
	return "12345"
}
