package models

//easyjson:json
type Thread struct {
	Id      int    `json:"id,omitempty"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum,omitempty"`
	Message string `json:"message"`
	Votes   int    `json:"votes,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Created string `json:"created,omitempty"`
}

type Vote struct {
	Nickname string `json:"nickname"`
	ThreadId int    `json:"thread_id"`
	Voice    int    `json:"voice"`
}

type Repository interface {
	CreateThread(thread Thread) (Thread, error)
	GetThreadBySlugOrId(slug string, id int) (Thread, error)
	GetThreads(slug string, limit int, since string, desc bool) ([]Thread, error)
	UpdateThreadBySlugOrId(thread Thread) (Thread, error)
	CreateVote(vote Vote) error
	UpdateVote(vote Vote) (int, error)
}

type Usecase interface {
	CreateThread(thread Thread) (Thread, error)
	GetThreadBySlugOrId(slugOrId string) (Thread, error)
	UpdateThreadBySlugOrId(thread Thread, slugOrId string) (Thread, error)
	GetThreads(slug string, limit int, since string, desc bool) ([]Thread, error)
	CreateVote(vote Vote, slugOrId string) (Thread, error)
}
