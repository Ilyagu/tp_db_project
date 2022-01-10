package repository

import (
	"context"
	forumModels "dbproject/internal/app/forum/models"
	userModels "dbproject/internal/app/user/models"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreForumRepo struct {
	Conn *pgxpool.Pool
}

func NewForumRepository(con *pgxpool.Pool) forumModels.Repository {
	return &PostgreForumRepo{con}
}

func (fr *PostgreForumRepo) CreateForum(forum forumModels.Forum) (forumModels.Forum, error) {
	var newForum forumModels.Forum
	err := fr.Conn.QueryRow(context.Background(), CreateForumQuery, forum.Title, forum.User, forum.Slug).
		Scan(&newForum.Title, &newForum.User, &newForum.Slug, &newForum.Posts, &newForum.Threads)
	if err != nil {
		return forumModels.Forum{}, err
	}
	return newForum, nil
}

func (fr *PostgreForumRepo) GetForum(slug string) (forumModels.Forum, error) {
	var forum forumModels.Forum
	err := fr.Conn.QueryRow(context.Background(), GetForumQuery, slug).
		Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return forumModels.Forum{}, err
	}
	return forum, nil
}

func (fr *PostgreForumRepo) GetForumUsers(forumSlug, since string, limit int, desc bool) ([]userModels.User, error) {
	forumUsersQuery := fmt.Sprintf(`select us.nickname, us.fullname, us.about, us.email from users us
							join users_to_forums utf on us.nickname = utf.nickname
						where utf.forum = '%s'`, forumSlug)
	if desc && since != "" {
		forumUsersQuery += fmt.Sprintf(` and utf.nickname < '%s'`, since)
	} else if since != "" {
		forumUsersQuery += fmt.Sprintf(` and utf.nickname > '%s'`, since)
	}
	forumUsersQuery += ` order by utf.nickname `
	if desc {
		forumUsersQuery += "desc"
	}
	forumUsersQuery += fmt.Sprintf(` limit %d`, limit)

	rows, err := fr.Conn.Query(context.Background(), forumUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	forumUsers := make([]userModels.User, 0)

	for rows.Next() {
		var forumUser userModels.User

		err := rows.Scan(&forumUser.Nickname, &forumUser.Fullname, &forumUser.About, &forumUser.Email)
		if err != nil {
			return nil, err
		}
		forumUsers = append(forumUsers, forumUser)
	}
	return forumUsers, nil
}
