package model

type User struct {
	ID       int
	Username string
	Email    string
	Status   bool
}

var users []User

func init() {

}

func NewUser(id int, username, email string) *User {
	return &User{
		ID:       id,
		Username: username,
		Email:    email,
	}
}

func GetUsers() []User {

	var users []User
	DB.Find(&users)
	return users
}

func GetUserCount() string {
	return "210"
}

func GetPageView() string {
	return "12345"
}
