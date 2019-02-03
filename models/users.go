package models

type User struct {
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserService interface {
	User(id int) (*User, error)
	Users() ([]*User, error)
	CreateUser(u *User) error
	DeleteUser(id int) error
	FindUser(email string) error
}

//Albums is an array of Album
type Users []User
