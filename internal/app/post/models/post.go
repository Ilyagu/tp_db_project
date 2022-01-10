package models

import (
	forumModels "dbproject/internal/app/forum/models"
	threadModels "dbproject/internal/app/thread/models"
	userModels "dbproject/internal/app/user/models"

	"github.com/jackc/pgtype"
)

//easyjson:json
type Post struct {
	Id       int              `json:"id,omitempty"`
	Parent   int              `json:"parent,omitempty"`
	Author   string           `json:"author"`
	Message  string           `json:"message"`
	IsEdited bool             `json:"isEdited,omitempty"`
	Forum    string           `json:"forum,omitempty"`
	Thread   int              `json:"thread,omitempty"`
	Created  string           `json:"created,omitempty"`
	Path     pgtype.Int8Array `json:"-"`
}

type PostUpdate struct {
	Message string `json:"message,omitempty"`
}

type PostFull struct {
	Post   Post                `json:"post,omitempty"`
	Author userModels.User     `json:"author,omitempty"`
	Thread threadModels.Thread `json:"thread,omitempty"`
	Forum  forumModels.Forum   `json:"forum,omitempty"`
}

type Status struct {
	User   int `json:"user"`
	Forum  int `json:"forum"`
	Thread int `json:"thread"`
	Post   int `json:"post"`
}

type Repository interface {
	CreatePosts(posts []Post, forumSlug string, threadId int) ([]Post, error)
	GetPostsFlat(threadId int, limit int, since int, desc bool) ([]Post, error)
	GetPostsTree(threadId int, limit int, since int, desc bool) ([]Post, error)
	GetPostsParentTree(threadId int, limit int, since int, desc bool) ([]Post, error)
	GetPost(postId int) (Post, error)
	UpdatePost(postId int, message string) (Post, error)
	ServiceStatus() (Status, error)
	ClearAll() error
}

type Usecase interface {
	CreatePosts(posts []Post, forumSlug string, threadId int) ([]Post, error)
	GetPosts(slugOrId string, limit int, since int, sort string, desc bool) ([]Post, error)
	GetPost(postId int, relatedStrs []string) (map[string]interface{}, error)
	UpdatePost(postId int, message string) (Post, error)
	ServiceStatus() (Status, error)
	ClearAll() error
}
