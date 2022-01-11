package usecase

import (
	forumModels "dbproject/internal/app/forum/models"
	"dbproject/internal/app/thread/models"
	threadModels "dbproject/internal/app/thread/models"
	userModels "dbproject/internal/app/user/models"
	"log"
	"strconv"

	"github.com/jackc/pgconn"
)

type ThreadUsecase struct {
	threadRepo threadModels.Repository
	userRepo   userModels.Repository
	forumRepo  forumModels.Repository
}

func NewThreadUsecase(tr threadModels.Repository, ur userModels.Repository, fr forumModels.Repository) threadModels.Usecase {
	return &ThreadUsecase{
		threadRepo: tr,
		userRepo:   ur,
		forumRepo:  fr,
	}
}

func (tu *ThreadUsecase) CreateThread(thread threadModels.Thread) (threadModels.Thread, error) {
	newThread, err := tu.threadRepo.CreateThread(thread)
	return newThread, err
}

func (tu *ThreadUsecase) GetThreadBySlugOrId(slugOrId string) (threadModels.Thread, error) {
	id, ok := strconv.Atoi(slugOrId)
	if ok != nil {
		id = 0
	}

	thread, err := tu.threadRepo.GetThreadBySlugOrId(slugOrId, id)
	if err != nil {
		return threadModels.Thread{}, err
	}

	return thread, err
}

func (tu *ThreadUsecase) UpdateThreadBySlugOrId(thread threadModels.Thread, slugOrId string) (threadModels.Thread, error) {
	id, ok := strconv.Atoi(slugOrId)
	if ok != nil {
		thread.Id = 0
	} else {
		thread.Id = id
	}
	thread.Slug = slugOrId
	existsThread, err := tu.threadRepo.GetThreadBySlugOrId(thread.Slug, thread.Id)
	if err != nil {
		log.Println(err)
		return threadModels.Thread{}, err
	}
	if thread.Title == "" {
		thread.Title = existsThread.Title
	}
	if thread.Message == "" {
		thread.Message = existsThread.Message
	}
	updatedThread, err := tu.threadRepo.UpdateThreadBySlugOrId(thread)
	if err != nil {
		return threadModels.Thread{}, err
	}

	return updatedThread, nil
}

func (tu *ThreadUsecase) GetThreads(slug string, limit int, since string, desc bool) ([]threadModels.Thread, error) {
	threads, err := tu.threadRepo.GetThreads(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}

	return threads, err
}

func (tu *ThreadUsecase) CreateVote(vote models.Vote, slugOrId string) (threadModels.Thread, error) {
	threadId, ok := strconv.Atoi(slugOrId)
	if ok != nil {
		threadId = 0
	}

	thread, err := tu.threadRepo.GetThreadBySlugOrId(slugOrId, threadId)
	if err != nil {
		return threadModels.Thread{}, err
	}

	_, err = tu.userRepo.GetUserByNickname(vote.Nickname)
	if err != nil {
		return threadModels.Thread{}, err
	}

	vote.ThreadId = thread.Id
	err = tu.threadRepo.CreateVote(vote)
	if err != nil {
		if err.(*pgconn.PgError).Code == "23503" {
			return threadModels.Thread{}, err
		}
		updatedVote, err := tu.threadRepo.UpdateVote(vote)
		if err != nil {
			return threadModels.Thread{}, err
		}
		if updatedVote != 0 {
			thread.Votes += 2 * vote.Voice
		}
		return thread, err
	}
	thread.Votes += vote.Voice
	return thread, err
}
