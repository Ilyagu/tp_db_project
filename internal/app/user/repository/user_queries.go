package repository

const (
	GetUserByNicknameQuery = `select nickname, fullname, about, email from users where nickname=$1;`

	GetUserByEmailQuery = `select nickname, fullname, about, email from users where email=$1;`

	GetUserByNicknameOrEmailQuery = `select nickname, fullname, about, email from users where nickname = $1 or email = $2`

	CreateUserQuery = `insert into users(nickname, fullname, about, email)
	values ($1,$2,$3,$4)
	returning nickname, fullname, about, email;`

	UpdateUserQuery = `update users set
		fullname=$2,
		about=$3,
		email=$4
	where nickname=$1
	returning nickname, fullname, about, email;`
)
