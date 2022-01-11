package delivery

import (
	forumModels "dbproject/internal/app/forum/models"
	"dbproject/internal/app/middlware"
	threadModels "dbproject/internal/app/thread/models"
	"dbproject/internal/pkg/responses"
	"dbproject/internal/pkg/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type ThreadHandler struct {
	threadUsecase threadModels.Usecase
	forumsUseacse forumModels.Usecase
}

func NewThreadHandler(router *router.Router, tu threadModels.Usecase, fu forumModels.Usecase) *ThreadHandler {
	threadHandler := &ThreadHandler{
		threadUsecase: tu,
		forumsUseacse: fu,
	}

	router.POST("/api/forum/{slug}/create", middlware.ReponseMiddlwareAndLogger(threadHandler.CreateThreadHandler))
	router.GET("/api/forum/{slug}/threads", middlware.ReponseMiddlwareAndLogger(threadHandler.GetThreads))
	router.GET("/api/thread/{slug_or_id}/details", middlware.ReponseMiddlwareAndLogger(threadHandler.GetThreadHandler))
	router.POST("/api/thread/{slug_or_id}/details", middlware.ReponseMiddlwareAndLogger(threadHandler.UpdateThreadHandler))
	router.POST("/api/thread/{slug_or_id}/vote", middlware.ReponseMiddlwareAndLogger(threadHandler.CreateVoteHandler))
	return threadHandler
}

func (th *ThreadHandler) CreateThreadHandler(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	newThread := threadModels.Thread{}
	err := json.Unmarshal(ctx.PostBody(), &newThread)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	newThread.Forum = slug

	if newThread.Slug != "" {
		exsistsThread, err := th.threadUsecase.GetThreadBySlugOrId(newThread.Slug)
		if err == nil {
			exsistsThreadBody, err := exsistsThread.MarshalJSON()
			if err != nil {
				log.Println(err)
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}
			ctx.SetStatusCode(http.StatusConflict)
			ctx.SetBody(exsistsThreadBody)
			return
		}
	}
	newThreadResp, err := th.threadUsecase.CreateThread(newThread)
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	if newThread.Slug == "" {
		newThreadResp.Slug = ""
	}
	newThreadRespBody, err := newThreadResp.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(newThreadRespBody)
}

func (th *ThreadHandler) GetThreadHandler(ctx *fasthttp.RequestCtx) {
	slug_or_id := ctx.UserValue("slug_or_id").(string)

	thread, err := th.threadUsecase.GetThreadBySlugOrId(slug_or_id)
	if err != nil {
		log.Println(err)
		responses.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	threadBody, err := thread.MarshalJSON()
	if err != nil {
		log.Println()
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(threadBody)
}

func (th *ThreadHandler) UpdateThreadHandler(ctx *fasthttp.RequestCtx) {
	slugOrId := ctx.UserValue("slug_or_id").(string)

	threadToUpdate := threadModels.Thread{}
	err := json.Unmarshal(ctx.PostBody(), &threadToUpdate)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	updatedThread, err := th.threadUsecase.UpdateThreadBySlugOrId(threadToUpdate, slugOrId)
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	updatedThreadBody, err := updatedThread.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(updatedThreadBody)
}

func (th *ThreadHandler) GetThreads(ctx *fasthttp.RequestCtx) {
	forumSlug, ok := ctx.UserValue("slug").(string)
	if !ok {
		responses.SendErrorResponse(ctx, http.StatusBadRequest, "bad request")
		return
	}

	limit, err := utils.ExtractIntValue(ctx, "limit")
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if limit == 0 {
		limit = 100
	}

	since := string(ctx.QueryArgs().Peek("since"))

	desc, err := utils.ExtractBoolValue(ctx, "desc")
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = th.forumsUseacse.GetForum(forumSlug)
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	threads, err := th.threadUsecase.GetThreads(forumSlug, limit, since, desc)
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	threadsBody, err := json.Marshal(threads)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(threadsBody)
}

func (th *ThreadHandler) CreateVoteHandler(ctx *fasthttp.RequestCtx) {
	slugOrId := ctx.UserValue("slug_or_id").(string)

	newVote := threadModels.Vote{}
	err := json.Unmarshal(ctx.PostBody(), &newVote)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	thread, err := th.threadUsecase.CreateVote(newVote, slugOrId)
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	threadBody, err := thread.MarshalJSON()
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(threadBody)
}
