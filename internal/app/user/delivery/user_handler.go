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

type Userandler struct {
	forumUsecase forumModels.Usecase
	userUsecase  userModels.Usecase
}

func NewUserHandler(router *router.Router, uu userModels.Usecase) {
	userHandler := &Userandler{
		userUsecase: uu,
	}

	router.POST("/api/user/{nickname}/create",
		middlware.ReponseMiddlwareAndLogger(userHandler.CreateUserHandler))
	router.GET("/api/user/{nickname}/profile",
		middlware.ReponseMiddlwareAndLogger(userHandler.GetUserHandler))
	router.POST("/api/user/{nickname}/profile",
		middlware.ReponseMiddlwareAndLogger(userHandler.UpdateUserHandler))
}

func (uh *Userandler) CreateUserHandler(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	newUser := userModels.User{}
	err := json.Unmarshal(ctx.PostBody(), &newUser)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	newUser.Nickname = nickname

	exsistsUsers, err := uh.userUsecase.GetUserByNicknameOrEmail(nickname, newUser.Email)
	if err != nil || len(exsistsUsers) == 0 {
		user, err := uh.userUsecase.CreateUser(newUser)
		if err != nil {
			log.Println(err)
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		userBody, err := user.MarshalJSON()
		if err != nil {
			log.Println(err)
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetStatusCode(http.StatusCreated)
		ctx.SetBody(userBody)
		return
	}
	existsUsersBody, err := json.Marshal(exsistsUsers)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusConflict)
	ctx.SetBody(existsUsersBody)
}

func (uh *Userandler) GetUserHandler(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	user, err := uh.userUsecase.GetUserByNickname(nickname)
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	userBody, err := user.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(userBody)
}

func (uh *Userandler) UpdateUserHandler(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	newUser := userModels.User{}
	err := json.Unmarshal(ctx.PostBody(), &newUser)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	newUser.Nickname = nickname

	existsUser, err := uh.userUsecase.GetUserByNickname(nickname)
	if err != nil {
		responses.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
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
		responses.SendErrorResponse(ctx, http.StatusConflict, err.Error())
		return
	}
	userUpdateBody, err := updateUser.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(userUpdateBody)
}
