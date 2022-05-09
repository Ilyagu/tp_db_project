package delivery

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	forumModels "dbproject/internal/app/forum/models"
	userModels "dbproject/internal/app/user/models"
	"dbproject/internal/pkg/responses"
)

type Userandler struct {
	forumUsecase forumModels.Usecase
	userUsecase  userModels.Usecase
}

func NewUserHandler(router *mux.Router, uu userModels.Usecase) {
	userHandler := &Userandler{
		userUsecase: uu,
	}

	router.HandleFunc("/api/user/{nickname}/create", userHandler.CreateUserHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/{nickname}/profile", userHandler.GetUserHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{nickname}/profile", userHandler.UpdateUserHandler).Methods("POST", "OPTIONS")
}

func (uh *Userandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]

	var newUser userModels.User
	err := responses.ReadJSON(r, &newUser)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	newUser.Nickname = nickname

	exsistsUsers, err := uh.userUsecase.GetUserByNicknameOrEmail(nickname, newUser.Email)
	if err != nil || len(exsistsUsers) == 0 {
		user, err := uh.userUsecase.CreateUser(newUser)
		if err != nil {
			log.Println(err)
			responses.SendWithoutBody(w, http.StatusInternalServerError)
			return
		}
		responses.Send(w, http.StatusCreated, user)
		return
	}
	responses.SendArray(w, http.StatusConflict, exsistsUsers)
}

func (uh *Userandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]

	user, err := uh.userUsecase.GetUserByNickname(nickname)
	if err != nil {
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}
	responses.Send(w, http.StatusOK, user)
}

func (uh *Userandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]

	var newUser userModels.User
	err := responses.ReadJSON(r, &newUser)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	newUser.Nickname = nickname

	existsUser, err := uh.userUsecase.GetUserByNickname(nickname)
	if err != nil {
		responses.SendError(w, http.StatusNotFound, err.Error())
		return
	}
	if newUser.About == "" {
		newUser.About = existsUser.About
	}
	if newUser.Fullname == "" {
		newUser.Fullname = existsUser.Fullname
	}
	if newUser.Email == "" {
		newUser.Email = existsUser.Email
	}
	updateUser, err := uh.userUsecase.UpdateUser(newUser)
	if err != nil {
		responses.SendError(w, http.StatusConflict, err.Error())
		return
	}
	responses.Send(w, http.StatusOK, updateUser)
}
