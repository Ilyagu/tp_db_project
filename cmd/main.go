package main

import (
	"context"
	fd "dbproject/internal/app/forum/delivery"
	fr "dbproject/internal/app/forum/repository"
	fu "dbproject/internal/app/forum/usecase"
	pd "dbproject/internal/app/post/delivery"
	pr "dbproject/internal/app/post/repository"
	pu "dbproject/internal/app/post/usecase"
	td "dbproject/internal/app/thread/delivery"
	tr "dbproject/internal/app/thread/repository"
	tu "dbproject/internal/app/thread/usecase"
	ud "dbproject/internal/app/user/delivery"
	ur "dbproject/internal/app/user/repository"
	uu "dbproject/internal/app/user/usecase"
	"log"

	fasthttpprom "dbproject/internal/pkg/metrics"

	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
)

func main() {
	router := router.New()
	p := fasthttpprom.NewPrometheus("")
	p.Use(router)

	dbpool, err := pgxpool.Connect(context.Background(),
		"host=89.208.196.139 port=5432 user=ilyagu dbname=forum password=password sslmode=disable",
	)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	log.Println("Success connection")
	defer dbpool.Close()

	// repositories
	userRepo := ur.NewUserRepository(dbpool)
	forumRepo := fr.NewForumRepository(dbpool)
	threadRepo := tr.NewThreadRepository(dbpool)
	postRepo := pr.NewPostRepository(dbpool)

	// usecases
	userUC := uu.NewUserUsecase(userRepo)
	forumUC := fu.NewForumUsecase(forumRepo, userRepo)
	threadUC := tu.NewThreadUsecase(threadRepo, userRepo, forumRepo)
	postUC := pu.NewPostUsecase(threadRepo, postRepo, userRepo, forumRepo)

	// delivety
	fd.NewForumHandler(router, forumUC, userUC)
	ud.NewUserHandler(router, userUC)
	td.NewThreadHandler(router, threadUC, forumUC)
	pd.NewPostHandler(router, threadUC, postUC)

	err = fasthttp.ListenAndServe(":5000", p.Handler)
	log.Fatal(err)
}
