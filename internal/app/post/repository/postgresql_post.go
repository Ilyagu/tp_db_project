package repository

import (
	"context"
	"dbproject/internal/app/post/models"
	"dbproject/internal/pkg/responses"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgrePostRepo struct {
	Conn *pgxpool.Pool
}

func NewPostRepository(con *pgxpool.Pool) models.Repository {
	return &PostgrePostRepo{con}
}

func (pr *PostgrePostRepo) CreatePosts(posts []models.Post, forumSlug string, threadId int) ([]models.Post, error) {
	newPosts := make([]models.Post, 0, 0)
	if len(posts) == 0 {
		return newPosts, nil
	}

	CreatePostsQuery := `insert into posts(
		parent,
		author,
		message,
		forum,
		thread,
		created) values `

	var valuesNames []string
	var values []interface{}
	timeForAllPosts := time.Now()
	i := 1
	for _, post := range posts {
		if post.Parent != 0 {
			var parentThreadId int
			query := `select thread from posts where id=$1`
			pr.Conn.QueryRow(context.Background(), query, post.Parent).Scan(&parentThreadId)
			if parentThreadId != threadId {
				log.Println(post.Parent)
				return nil, errors.New("lolahahahahah")
			}
		}

		var existsUser bool
		query := "select exists(select nickname from users where nickname=$1)"
		pr.Conn.QueryRow(context.Background(), query, post.Author).Scan(&existsUser)
		if !existsUser {
			return nil, responses.UserNotExsists
		}

		valuesNames = append(valuesNames, fmt.Sprintf(
			"(nullif($%d, 0), $%d, $%d, $%d, $%d, $%d)",
			i, i+1, i+2, i+3, i+4, i+5))
		i += 6
		values = append(values, post.Parent, post.Author, post.Message, forumSlug, threadId, timeForAllPosts)
	}

	CreatePostsQuery += strings.Join(valuesNames[:], ",")
	CreatePostsQuery += `returning id, case when parent is null then 0 else parent end as parent,
	author, message, is_edited, forum, thread, created, path`

	rows, err := pr.Conn.Query(context.Background(), CreatePostsQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		createdTime := &time.Time{}

		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, createdTime, &post.Path)
		if err != nil {
			return nil, err
		}

		post.Created = strfmt.DateTime(createdTime.UTC()).String()
		newPosts = append(newPosts, post)
	}

	return newPosts, nil
}

func (pr *PostgrePostRepo) GetPostsFlat(threadId int, limit int, since int, desc bool) ([]models.Post, error) {

	query := `select id, case when parent is null then 0 else parent end as parent,
	author, message, is_edited, forum, thread, created, path from posts where thread=$1 `

	if desc {
		if since > 0 {
			query += fmt.Sprintf("and id < %d ", since)
		}
		query += `order by id desc `
	} else {
		if since > 0 {
			query += fmt.Sprintf("and id > %d ", since)
		}
		query += `order by id `
	}
	query += `limit nullif($2, 0)`
	var posts []models.Post

	rows, err := pr.Conn.Query(context.Background(), query, threadId, limit)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		createdTime := &time.Time{}

		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, createdTime, &post.Path)
		if err != nil {
			return posts, err
		}

		post.Created = strfmt.DateTime(createdTime.UTC()).String()
		posts = append(posts, post)

	}
	return posts, err
}

func (pr *PostgrePostRepo) GetPostsTree(threadId int, limit int, since int, desc bool) ([]models.Post, error) {
	var treeQuery string
	sinceQuery := ""
	if since != 0 {
		if desc {
			sinceQuery = `and path < `
		} else {
			sinceQuery = `and path > `
		}
		sinceQuery += fmt.Sprintf(`(select path from posts where id = %d)`, since)
	}
	if desc {
		treeQuery = fmt.Sprintf(
			`select id, case when parent is null then 0 else parent end as parent,
			author, message, is_edited, forum, thread, created, path
			from posts where thread=$1 %s order by path desc, id desc limit nullif($2, 0);`, sinceQuery)
	} else {
		treeQuery = fmt.Sprintf(
			`select id, case when parent is null then 0 else parent end as parent,
			author, message, is_edited, forum, thread, created, path
			from posts where thread=$1 %s order by path, id limit nullif($2, 0);`, sinceQuery)
	}

	var posts []models.Post

	rows, err := pr.Conn.Query(context.Background(), treeQuery, threadId, limit)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		createdTime := &time.Time{}

		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, createdTime, &post.Path)
		if err != nil {
			return posts, err
		}

		post.Created = strfmt.DateTime(createdTime.UTC()).String()
		posts = append(posts, post)

	}
	return posts, err
}

func (pr *PostgrePostRepo) GetPostsParentTree(threadId int, limit int, since int, desc bool) ([]models.Post, error) {
	var query string
	sinceQuery := ""
	if since != 0 {
		if desc {
			sinceQuery = `and path[1] < `
		} else {
			sinceQuery = `and path[1] > `
		}
		sinceQuery += fmt.Sprintf(`(select path[1] from posts where id = %d)`, since)
	}

	parentsQuery := fmt.Sprintf(
		`select id from posts where thread = $1 and parent is null %s`, sinceQuery)

	if desc {
		parentsQuery += `order by id desc`
		if limit > 0 {
			parentsQuery += fmt.Sprintf(` limit %d`, limit)
		}
		query = fmt.Sprintf(
			`select id, case when parent is null then 0 else parent end as parent,
			author, message, is_edited, forum, thread, created, path 
			from posts where path[1] in (%s) order by path[1] desc, path, id;`, parentsQuery)
	} else {
		parentsQuery += `order by id`
		if limit > 0 {
			parentsQuery += fmt.Sprintf(` limit %d`, limit)
		}
		query = fmt.Sprintf(
			`select id, case when parent is null then 0 else parent end as parent,
			author, message, is_edited, forum, thread, created, path 
			from posts where path[1] in (%s) order by path,id;`, parentsQuery)
	}
	var posts []models.Post
	rows, err := pr.Conn.Query(context.Background(), query, threadId)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		createdTime := &time.Time{}

		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, createdTime, &post.Path)
		if err != nil {
			return posts, err
		}

		post.Created = strfmt.DateTime(createdTime.UTC()).String()
		posts = append(posts, post)

	}
	return posts, err
}

func (pr *PostgrePostRepo) GetPost(postId int) (models.Post, error) {
	var post models.Post
	createdTime := &time.Time{}

	err := pr.Conn.QueryRow(context.Background(), GetPostQuery, postId).
		Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, createdTime, &post.Path)
	if err != nil {
		return models.Post{}, err
	}

	post.Created = strfmt.DateTime(createdTime.UTC()).String()

	return post, nil
}

func (pr *PostgrePostRepo) UpdatePost(postId int, message string) (models.Post, error) {
	var post models.Post
	createdTime := &time.Time{}

	err := pr.Conn.QueryRow(context.Background(), UpdatePostQuery, postId, message).
		Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited,
			&post.Forum, &post.Thread, createdTime, &post.Path)
	if err != nil {
		return models.Post{}, err
	}

	post.Created = strfmt.DateTime(createdTime.UTC()).String()

	return post, nil
}

func (pr *PostgrePostRepo) ServiceStatus() (models.Status, error) {
	var status models.Status

	err := pr.Conn.QueryRow(context.Background(), GetStatusQuery).
		Scan(&status.User, &status.Forum, &status.Thread, &status.Post)
	if err != nil {
		return models.Status{}, err
	}

	return status, nil
}

func (pr *PostgrePostRepo) ClearAll() error {
	_, err := pr.Conn.Exec(context.Background(), ClearAllQuery)
	if err != nil {
		return err
	}

	return nil
}
