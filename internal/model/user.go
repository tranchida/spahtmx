package model

import (
	"context"
	"time"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type Base struct {
 ID        string     `gorm:"type:uuid;primary_key;"`
 CreatedAt time.Time  `json:"created_at"`
 UpdatedAt time.Time  `json:"updated_at"`
 DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
 b.ID = ksuid.New().String()
 return
}


type User struct {
	Base
	Username string
	Email    string
	Status   bool
}

func GetUsers() []User {

	var users []User
	DB.Find(&users)
	return users
}

func UpdateUserStatus(ctx context.Context, id string) {

	genericDB := gorm.G[User](DB)
	user, err := genericDB.Where("id = ?", id).Take(ctx)
	if err != nil {
		return
	}
	_, err = genericDB.Where("id = ?", id).Update(ctx, "status", !user.Status)
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
