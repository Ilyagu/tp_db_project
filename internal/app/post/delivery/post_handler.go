package delivery

import (
	"dbproject/internal/app/middlware"
	postModels "dbproject/internal/app/post/models"
	threadModels "dbproject/internal/app/thread/models"
	"dbproject/internal/pkg/responses"
	"dbproject/internal/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type PostHandler struct {
	threadUsecase threadModels.Usecase
	postUsecase   postModels.Usecase
}

func NewPostHandler(router *mux.Router, tu threadModels.Usecase, pu postModels.Usecase) *PostHandler {
	postHandler := &PostHandler{
		threadUsecase: tu,
		postUsecase:   pu,
	}

	router.POST("/api/thread/{slug_or_id}/create", middlware.ReponseMiddlwareAndLogger(postHandler.CreatePostHandler))
	router.GET("/api/thread/{slug_or_id}/posts", middlware.ReponseMiddlwareAndLogger(postHandler.GetPostsHandler))
	router.GET("/api/post/{id}/details", middlware.ReponseMiddlwareAndLogger(postHandler.GetPostFullHandler))
	router.POST("/api/post/{id}/details", middlware.ReponseMiddlwareAndLogger(postHandler.UpdatePostHandler))
	router.POST("/api/service/clear", middlware.ReponseMiddlwareAndLogger(postHandler.ClearAllHandler))
	router.GET("/api/service/status", middlware.ReponseMiddlwareAndLogger(postHandler.ServiceStatusHandler))

	return postHandler
}

func (ph *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]

	var newPosts []postModels.Post
	err := responses.ReadJSONArray(r, &newPosts)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	existsThread, err := ph.threadUsecase.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	if len(newPosts) == 0 {
		responses.SendArray(w, http.StatusCreated, newPosts)
		return
	}

	newPostsResp, err := ph.postUsecase.CreatePosts(newPosts, existsThread.Forum, existsThread.Id)
	if err == responses.UserNotExsists {
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		responses.SendError(w, http.StatusConflict, err.Error())
		return
	}
	responses.SendArray(w, http.StatusCreated, newPostsResp)
}

func (ph *PostHandler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	threadSlugOrID := mux.Vars(r)["slug_or_id"]

	limit, err := utils.ExtractIntValue(r, "limit")
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	since, err := utils.ExtractIntValue(r, "since")
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sortType := r.URL.Query().Get("sort")
	if sortType == "" {
		sortType = "flat"
	}

	desc, err := utils.ExtractBoolValue(r, "desc")
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	posts, err := ph.postUsecase.GetPosts(threadSlugOrID, limit, since, sortType, desc)
	if err != nil {
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	if posts == nil {
		nullPosts := make([]postModels.Post, 0)
		nullPostsBody, err := json.Marshal(nullPosts)
		if err != nil {
			responses.SendError(w, http.StatusInternalServerError, err.Error())
			return
		}
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(nullPostsBody)
		return
	}
	postsBody, err := json.Marshal(posts)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(postsBody)
	return
}

func (ph *PostHandler) GetPostFullHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	relatedArr := string(ctx.QueryArgs().Peek("related"))
	relatedStrs := strings.Split(relatedArr, ",")
	for len(relatedStrs) < 3 {
		relatedStrs = append(relatedStrs, "")
	}
	if err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	post, err := ph.postUsecase.GetPost(postId, relatedStrs)
	if err != nil {
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	postBody, err := json.Marshal(post)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(postBody)
}

func (ph *PostHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	updateMessage := postModels.PostUpdate{}
	err = json.Unmarshal(ctx.PostBody(), &updateMessage)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	post, err := ph.postUsecase.UpdatePost(postId, updateMessage.Message)
	if err != nil {
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	postBody, err := json.Marshal(post)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(postBody)
}

func (ph *PostHandler) ServiceStatusHandler(w http.ResponseWriter, r *http.Request) {
	status, err := ph.postUsecase.ServiceStatus()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	statusBody, err := status.MarshalJSON()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(statusBody)
}

func (ph *PostHandler) ClearAllHandler(w http.ResponseWriter, r *http.Request) {
	err := ph.postUsecase.ClearAll()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.SetStatusCode(http.StatusOK)
}
