package model

type User struct {
	ID       int
	Username string
	Email    string
	Status   bool
}

var users []User

func init() {
	// Initialisation si nécessaire
	users = []User{
		{ID: 1, Username: "alice", Email: "alice@fake.com", Status: true},
		{ID: 2, Username: "bob", Email: "bob@fake.com", Status: false},
		{ID: 3, Username: "charlie", Email: "charlie@fake.com", Status: true},
	}

}

func NewUser(id int, username, email string) *User {
	return &User{
		ID:       id,
		Username: username,
		Email:    email,
	}
}

func GetUsers() []User {
	return users
}
