package usecase

import (
	forumModels "dbproject/internal/app/forum/models"
	"dbproject/internal/app/post/models"
	postModels "dbproject/internal/app/post/models"
	threadModels "dbproject/internal/app/thread/models"
	userModels "dbproject/internal/app/user/models"
	"dbproject/internal/pkg/utils"
	"errors"
	"strconv"
)

type PostUsecase struct {
	threadRepo threadModels.Repository
	postRepo   postModels.Repository
	userRepo   userModels.Repository
	forumRepo  forumModels.Repository
}

func NewPostUsecase(tr threadModels.Repository, pr postModels.Repository, ur userModels.Repository, fr forumModels.Repository) postModels.Usecase {
	return &PostUsecase{
		threadRepo: tr,
		postRepo:   pr,
		userRepo:   ur,
		forumRepo:  fr,
	}
}

func (pu *PostUsecase) CreatePosts(posts []models.Post, forumSlug string, threadId int) ([]models.Post, error) {
	newPosts, err := pu.postRepo.CreatePosts(posts, forumSlug, threadId)
	return newPosts, err
}

func (pu *PostUsecase) GetPosts(slugOrId string, limit int, since int, sort string, desc bool) ([]models.Post, error) {
	var err error
	threadId, ok := strconv.Atoi(slugOrId)
	if ok != nil {
		threadId = 0
	}

	thread, err := pu.threadRepo.GetThreadBySlugOrId(slugOrId, threadId)
	if err != nil {
		return nil, err
	}

	switch sort {
	case "flat":
		return pu.postRepo.GetPostsFlat(thread.Id, limit, since, desc)
	case "tree":
		return pu.postRepo.GetPostsTree(thread.Id, limit, since, desc)
	case "parent_tree":
		return pu.postRepo.GetPostsParentTree(thread.Id, limit, since, desc)
	default:
		return nil, errors.New("THERE IS NO SORT WITH THIS NAME")
	}
}

func (pu *PostUsecase) GetPost(postId int, relatedStrs []string) (map[string]interface{}, error) {
	post, err := pu.postRepo.GetPost(postId)
	if err != nil {
		return nil, err
	}

	postFullMap := map[string]interface{}{
		"post": post,
	}

	if utils.Find(relatedStrs, "user") {
		user, err := pu.userRepo.GetUserByNickname(post.Author)
		if err != nil {
			return nil, err
		}
		postFullMap["author"] = user
	}
	if utils.Find(relatedStrs, "forum") {
		forum, err := pu.forumRepo.GetForum(post.Forum)
		if err != nil {
			return nil, err
		}
		postFullMap["forum"] = forum
	}
	if utils.Find(relatedStrs, "thread") {
		thread, err := pu.threadRepo.GetThreadById(post.Thread)
		if err != nil {
			return nil, err
		}
		postFullMap["thread"] = thread
	}

	return postFullMap, err
}

func (pu *PostUsecase) UpdatePost(postId int, message string) (models.Post, error) {
	post, err := pu.postRepo.UpdatePost(postId, message)
	return post, err
}

func (pu *PostUsecase) ServiceStatus() (models.Status, error) {
	status, err := pu.postRepo.ServiceStatus()
	return status, err
}

func (pu *PostUsecase) ClearAll() error {
	err := pu.postRepo.ClearAll()
	return err
}
