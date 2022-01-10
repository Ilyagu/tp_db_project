package models

import "dbproject/internal/app/user/models"

//easyjson:json
type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int    `json:"posts,omitempty"`
	Threads int    `json:"threads,omitempty"`
}

type Repository interface {
	CreateForum(forum Forum) (Forum, error)
	GetForum(slug string) (Forum, error)
	GetForumUsers(forumSlug, since string, limit int, desc bool) ([]models.User, error)
}

type Usecase interface {
	CreateForum(forum Forum) (Forum, error)
	GetForum(slug string) (Forum, error)
	GetForumUsers(forumSlug, since string, limit int, desc bool) ([]models.User, error)
}
