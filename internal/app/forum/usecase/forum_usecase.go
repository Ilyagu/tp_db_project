package usecase

import (
	"dbproject/internal/app/forum/models"
	forumModels "dbproject/internal/app/forum/models"
	userModels "dbproject/internal/app/user/models"
)

type ForumUsecase struct {
	forumRepo forumModels.Repository
	userRepo  userModels.Repository
}

func NewForumUsecase(fr forumModels.Repository, ur userModels.Repository) models.Usecase {
	return &ForumUsecase{
		forumRepo: fr,
		userRepo:  ur,
	}
}

func (fu *ForumUsecase) CreateForum(forum forumModels.Forum) (forumModels.Forum, error) {
	newForum, err := fu.forumRepo.CreateForum(forum)
	user, err := fu.userRepo.GetUserByNickname(newForum.User)
	newForum.User = user.Nickname
	return newForum, err
}

func (fu *ForumUsecase) GetForum(slug string) (forumModels.Forum, error) {
	forum, err := fu.forumRepo.GetForum(slug)
	user, err := fu.userRepo.GetUserByNickname(forum.User)
	forum.User = user.Nickname
	return forum, err
}

func (fu *ForumUsecase) GetForumUsers(forumSlug, since string, limit int, desc bool) ([]userModels.User, error) {
	_, err := fu.forumRepo.GetForum(forumSlug)
	if err != nil {
		return nil, err
	}
	forumUsers, err := fu.forumRepo.GetForumUsers(forumSlug, since, limit, desc)
	if err != nil {
		return nil, err
	}

	return forumUsers, nil
}
