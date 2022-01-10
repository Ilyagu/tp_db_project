package repository

const (
	CreateForumQuery = `insert into forums(title, "user", slug) values ($1, $2, $3)
	returning title, "user", slug, posts, threads;`

	GetForumQuery = `select title, "user", slug, posts, threads from forums where slug=$1;`
)
