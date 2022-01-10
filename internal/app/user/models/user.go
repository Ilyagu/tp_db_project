package models

//easyjson:json
type User struct {
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
}

type Repository interface {
	CreateUser(user User) (User, error)
	GetUserByNicknameOrEmail(nickname, email string) ([]User, error)
	GetUserByNickname(nickname string) (User, error)
	UpdateUser(user User) (User, error)
	GetUserByEmail(email string) (User, error)
}

type Usecase interface {
	CreateUser(user User) (User, error)
	GetUserByNicknameOrEmail(nickname, email string) ([]User, error)
	GetUserByNickname(nickname string) (User, error)
	UpdateUser(user User) (User, error)
	GetUserByEmail(email string) (User, error)
}
