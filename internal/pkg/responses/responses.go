package responses

import (
	"errors"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
)

var (
	UserNotExsists = errors.New("user not exists")
)

//easyjson:json
type Response struct {
	Message string `json:"message"`
}

func SendErrorResponse(ctx *fasthttp.RequestCtx, status int, message string) {
	log.Println(message)
	resp := &Response{
		Message: message,
	}
	respBody, err := resp.MarshalJSON()
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(status)
	ctx.SetBody(respBody)
}
