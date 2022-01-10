package repository

import (
	"context"
	"dbproject/internal/app/user/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreUserRepo struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(con *pgxpool.Pool) models.Repository {
	return &PostgreUserRepo{con}
}

func (ur *PostgreUserRepo) GetUserByNickname(nickname string) (models.User, error) {
	var user models.User
	err := ur.Conn.QueryRow(context.Background(), GetUserByNicknameQuery, nickname).
		Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (ur *PostgreUserRepo) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := ur.Conn.QueryRow(context.Background(), GetUserByEmailQuery, email).
		Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (userRep *PostgreUserRepo) GetUserByNicknameOrEmail(nickname, email string) ([]models.User, error) {
	users := make([]models.User, 0)

	rows, err := userRep.Conn.Query(context.Background(), GetUserByNicknameOrEmailQuery, nickname, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		users = append(users, user)
	}
	return users, nil
}

func (ur *PostgreUserRepo) CreateUser(user models.User) (models.User, error) {
	var newUser models.User
	err := ur.Conn.QueryRow(context.Background(), CreateUserQuery,
		user.Nickname, user.Fullname, user.About, user.Email).
		Scan(&newUser.Nickname, &newUser.Fullname, &newUser.About, &newUser.Email)
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

func (ur *PostgreUserRepo) UpdateUser(user models.User) (models.User, error) {
	var newUser models.User
	err := ur.Conn.QueryRow(context.Background(), UpdateUserQuery,
		user.Nickname, user.Fullname, user.About, user.Email).
		Scan(&newUser.Nickname, &newUser.Fullname, &newUser.About, &newUser.Email)
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}
