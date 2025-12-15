package domain

import (
	"time"
)

type User struct {
	ID        string
	Username  string
	Email     string
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
