package delivery

import (
	forumModels "dbproject/internal/app/forum/models"
	threadModels "dbproject/internal/app/thread/models"
	"dbproject/internal/pkg/responses"
	"dbproject/internal/pkg/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ThreadHandler struct {
	threadUsecase threadModels.Usecase
	forumsUseacse forumModels.Usecase
}

func NewThreadHandler(router *mux.Router, tu threadModels.Usecase, fu forumModels.Usecase) *ThreadHandler {
	threadHandler := &ThreadHandler{
		threadUsecase: tu,
		forumsUseacse: fu,
	}

	router.HandleFunc("/api/forum/{slug}/create", threadHandler.CreateThreadHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/forum/{slug}/threads", threadHandler.GetThreads).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/thread/{slug_or_id}/details", threadHandler.GetThreadHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/thread/{slug_or_id}/details", threadHandler.UpdateThreadHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/thread/{slug_or_id}/vote", threadHandler.CreateVoteHandler).Methods("POST", "OPTIONS")
	return threadHandler
}

func (th *ThreadHandler) CreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	var newThread threadModels.Thread
	err := responses.ReadJSON(r, &newThread)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	newThread.Forum = slug

	if newThread.Slug != "" {
		exsistsThread, err := th.threadUsecase.GetThreadBySlugOrId(newThread.Slug)
		if err == nil {
			responses.Send(w, http.StatusConflict, exsistsThread)
			return
		}
	}
	newThreadResp, err := th.threadUsecase.CreateThread(newThread)
	if err != nil {
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}
	if newThread.Slug == "" {
		newThreadResp.Slug = ""
	}
	responses.Send(w, http.StatusCreated, newThreadResp)
}

func (th *ThreadHandler) GetThreadHandler(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]

	thread, err := th.threadUsecase.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		log.Println(err)
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	responses.Send(w, http.StatusOK, thread)
}

func (th *ThreadHandler) UpdateThreadHandler(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]

	var threadToUpdate threadModels.Thread
	err := responses.ReadJSON(r, &threadToUpdate)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedThread, err := th.threadUsecase.UpdateThreadBySlugOrId(threadToUpdate, slugOrId)
	if err != nil {
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}
	responses.Send(w, http.StatusOK, updatedThread)
}

func (th *ThreadHandler) GetThreads(w http.ResponseWriter, r *http.Request) {
	forumSlug := mux.Vars(r)["slug_or_id"]

	limit, err := utils.ExtractIntValue(r, "limit")
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if limit == 0 {
		limit = 100
	}

	since := r.URL.Query().Get("since")

	desc, err := utils.ExtractBoolValue(r, "desc")
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = th.forumsUseacse.GetForum(forumSlug)
	if err != nil {
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}
	threads, err := th.threadUsecase.GetThreads(forumSlug, limit, since, desc)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responses.SendArray(w, http.StatusOK, threads)
}

func (th *ThreadHandler) CreateVoteHandler(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]

	var newVote threadModels.Vote
	err := responses.ReadJSON(r, &newVote)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	thread, err := th.threadUsecase.CreateVote(newVote, slugOrId)
	if err != nil {
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}
	responses.Send(w, http.StatusOK, thread)
}
