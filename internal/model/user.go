package model

type User struct {
	ID       int
	Username string
	Email    string
	Status   bool
}

var users []User

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

func UpdateUserStatus(id int) {
	var user User
	DB.First(&user, id)
	user.Status = !user.Status
	DB.Save(&user)
}

func GetUserCount() string {
	return "210"
}

func GetPageView() string {
	return "12345"
}
