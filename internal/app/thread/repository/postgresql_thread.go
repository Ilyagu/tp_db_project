package repository

import (
	"context"
	"dbproject/internal/app/thread/models"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreThreadRepo struct {
	Conn *pgxpool.Pool
}

func NewThreadRepository(con *pgxpool.Pool) models.Repository {
	return &PostgreThreadRepo{con}
}

func (tr *PostgreThreadRepo) CreateThread(thread models.Thread) (models.Thread, error) {
	var newThread models.Thread
	var err error

	createdTime := &time.Time{}
	if thread.Created == "" {
		err = tr.Conn.QueryRow(context.Background(), CreateThreadNoCreatedQuery,
			thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug).
			Scan(&newThread.Id, &newThread.Title, &newThread.Author,
				&newThread.Forum, &newThread.Message, &newThread.Slug,
				&newThread.Votes, createdTime)
	} else {
		err = tr.Conn.QueryRow(context.Background(), CreateThreadQuery,
			thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created).
			Scan(&newThread.Id, &newThread.Title, &newThread.Author,
				&newThread.Forum, &newThread.Message, &newThread.Slug,
				&newThread.Votes, createdTime)
	}
	if err != nil {
		return models.Thread{}, err
	}
	newThread.Created = strfmt.DateTime(createdTime.UTC()).String()

	return newThread, nil
}

func (tr *PostgreThreadRepo) GetThreadBySlugOrId(slug string, id int) (models.Thread, error) {
	var exsistsThread models.Thread
	createdTime := &time.Time{}
	err := tr.Conn.QueryRow(context.Background(), GetThreadBySlugOrIdQuery, id, slug).
		Scan(&exsistsThread.Id, &exsistsThread.Title, &exsistsThread.Author,
			&exsistsThread.Forum, &exsistsThread.Message, &exsistsThread.Slug,
			&exsistsThread.Votes, createdTime)
	if err != nil {
		return models.Thread{}, err
	}
	exsistsThread.Created = strfmt.DateTime(createdTime.UTC()).String()

	return exsistsThread, nil
}

func (tr *PostgreThreadRepo) GetThreadById(id int) (models.Thread, error) {
	var exsistsThread models.Thread
	createdTime := &time.Time{}
	err := tr.Conn.QueryRow(context.Background(), GetThreadByIdQuery, id).
		Scan(&exsistsThread.Id, &exsistsThread.Title, &exsistsThread.Author,
			&exsistsThread.Forum, &exsistsThread.Message, &exsistsThread.Slug,
			&exsistsThread.Votes, createdTime)
	if err != nil {
		return models.Thread{}, err
	}
	exsistsThread.Created = strfmt.DateTime(createdTime.UTC()).String()

	return exsistsThread, nil
}

func (tr *PostgreThreadRepo) UpdateThreadBySlugOrId(thread models.Thread) (models.Thread, error) {
	var updatedThread models.Thread
	createdTime := &time.Time{}
	err := tr.Conn.QueryRow(context.Background(), UpdateThreadQuery,
		thread.Title, thread.Message, thread.Slug, thread.Id).
		Scan(&updatedThread.Id, &updatedThread.Title, &updatedThread.Author,
			&updatedThread.Forum, &updatedThread.Message, &updatedThread.Slug,
			&updatedThread.Votes, createdTime)
	if err != nil {
		return models.Thread{}, err
	}
	updatedThread.Created = strfmt.DateTime(createdTime.UTC()).String()

	return updatedThread, nil
}

func (tr *PostgreThreadRepo) GetThreads(slug string, limit int, since string, desc bool) ([]models.Thread, error) {
	var createdExpression string
	var orderExpression string

	if since != "" && desc {
		createdExpression = fmt.Sprintf("and created <= '%s'", since)
	} else if since != "" && !desc {
		createdExpression = fmt.Sprintf("and created >= '%s'", since)
	}

	if desc {
		orderExpression = "desc"
	}

	getThreadsQuery := fmt.Sprintf("select * from threads where forum=$1 %s order by created %s limit nullif($2, 0)",
		createdExpression, orderExpression)

	threads := make([]models.Thread, 0, 0)
	rows, err := tr.Conn.Query(context.Background(), getThreadsQuery, slug, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var exsistsThread models.Thread
		createdTime := &time.Time{}

		err = rows.Scan(&exsistsThread.Id, &exsistsThread.Title, &exsistsThread.Author,
			&exsistsThread.Forum, &exsistsThread.Message, &exsistsThread.Slug,
			&exsistsThread.Votes, createdTime)
		if err != nil {
			return nil, err
		}

		exsistsThread.Created = strfmt.DateTime(createdTime.UTC()).String()

		threads = append(threads, exsistsThread)
	}

	return threads, err
}

func (tr *PostgreThreadRepo) CreateVote(vote models.Vote) error {
	_, err := tr.Conn.Exec(context.Background(), CreateVoteQuery, vote.Nickname,
		vote.ThreadId, vote.Voice)
	return err
}

func (tr *PostgreThreadRepo) UpdateVote(vote models.Vote) (int, error) {
	res, err := tr.Conn.Exec(context.Background(), UpdateVoteQuery, vote.Voice, vote.ThreadId, vote.Nickname)
	return int(res.RowsAffected()), err
}
