package usecase

import (
	"dbproject/internal/app/user/models"
)

type UserUsecase struct {
	userRep models.Repository
}

func NewUserUsecase(ur models.Repository) models.Usecase {
	return &UserUsecase{
		userRep: ur,
	}
}

func (uu *UserUsecase) GetUserByNickname(nickname string) (models.User, error) {
	user, err := uu.userRep.GetUserByNickname(nickname)
	return user, err
}

func (uu *UserUsecase) GetUserByEmail(email string) (models.User, error) {
	user, err := uu.userRep.GetUserByEmail(email)
	return user, err
}

func (userUse *UserUsecase) GetUserByNicknameOrEmail(nickname, email string) ([]models.User, error) {
	users, err := userUse.userRep.GetUserByNicknameOrEmail(nickname, email)
	return users, err
}

func (uu *UserUsecase) CreateUser(user models.User) (models.User, error) {
	user, err := uu.userRep.CreateUser(user)
	return user, err
}

func (uu *UserUsecase) UpdateUser(user models.User) (models.User, error) {
	user, err := uu.userRep.UpdateUser(user)
	return user, err
}
