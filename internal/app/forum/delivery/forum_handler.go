package delivery

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	forumModels "dbproject/internal/app/forum/models"
	"dbproject/internal/app/middlware"
	userModels "dbproject/internal/app/user/models"
	"dbproject/internal/pkg/responses"
)

type ForumHandler struct {
	forumUsecase forumModels.Usecase
	userUsecase  userModels.Usecase
}

func NewForumHandler(router *router.Router, fu forumModels.Usecase, uu userModels.Usecase) {
	forumHandler := &ForumHandler{
		forumUsecase: fu,
		userUsecase:  uu,
	}

	router.POST("/api/forum/create",
		middlware.ReponseMiddlwareAndLogger(forumHandler.CreateForumHandler))
	router.GET("/api/forum/{slug}/details",
		middlware.ReponseMiddlwareAndLogger(forumHandler.GetForumHandler))
	router.GET("/api/forum/{forum_slug}/users",
		middlware.ReponseMiddlwareAndLogger(forumHandler.GetForumUsersHandler))
}

func (fh *ForumHandler) CreateForumHandler(ctx *fasthttp.RequestCtx) {
	newForum := forumModels.Forum{}
	err := json.Unmarshal(ctx.PostBody(), &newForum)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	_, err = fh.userUsecase.GetUserByNickname(newForum.User)
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	forum, err := fh.forumUsecase.GetForum(newForum.Slug)
	if err != nil {
		log.Println(err)
		newForumResp, err := fh.forumUsecase.CreateForum(newForum)
		if err != nil {
			log.Println(err)
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		newForumBody, err := newForumResp.MarshalJSON()
		if err != nil {
			log.Println(err)
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetStatusCode(http.StatusCreated)
		ctx.SetBody(newForumBody)
		return
	}
	ctx.SetStatusCode(http.StatusConflict)
	forumBody, err := forum.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetBody(forumBody)
}

func (fh *ForumHandler) GetForumHandler(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	forum, err := fh.forumUsecase.GetForum(slug)
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	forumBody, err := forum.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(forumBody)
}

func (fh *ForumHandler) GetForumUsersHandler(ctx *fasthttp.RequestCtx) {
	forumSlug := ctx.UserValue("forum_slug").(string)
	desc := ctx.QueryArgs().GetBool("desc")
	limit, err := ctx.QueryArgs().GetUint("limit")
	if err != nil {
		limit = 100
	}
	since := string(ctx.QueryArgs().Peek("since"))

	forumUsers, err := fh.forumUsecase.GetForumUsers(forumSlug, since, limit, desc)
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	forumUsersBody, err := json.Marshal(forumUsers)
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(forumUsersBody)
}
