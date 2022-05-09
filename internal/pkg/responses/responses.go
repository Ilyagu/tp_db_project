package responses

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
)

var (
	UserNotExsists = errors.New("user not exists")
)

//easyjson:json
type ErrorResponse struct {
	Message string `json:"message"`
}

type ReadModel interface {
	UnmarshalJSON(data []byte) error
}

type WriteModel interface {
	MarshalJSON() ([]byte, error)
}

func SendErrorResponse(ctx *fasthttp.RequestCtx, status int, message string) {
	log.Println(message)
	resp := &ErrorResponse{
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

func Send(w http.ResponseWriter, respCode int, respBody WriteModel) {
	w.WriteHeader(respCode)
	_ = WriteJSON(w, respBody)
}

func SendArray(w http.ResponseWriter, respCode int, respBody interface{}) {
	w.WriteHeader(respCode)
	_ = WriteJSONArray(w, respBody)
}

func SendError(w http.ResponseWriter, respCode int, errorMsg string) {
	Send(w, respCode, ErrorResponse{
		Message: errorMsg,
	})
}

func SendWithoutBody(w http.ResponseWriter, respCode int) {
	w.WriteHeader(respCode)
}

func ReadJSON(r *http.Request, data ReadModel) error {
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = data.UnmarshalJSON(byteReq)
	if err != nil {
		return err
	}

	return nil
}

func ReadJSONArray(r *http.Request, data interface{}) error {
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteReq, &data)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSONArray(w http.ResponseWriter, data interface{}) error {
	byteResp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(byteResp)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, data WriteModel) error {
	byteResp, err := data.MarshalJSON()
	if err != nil {
		return err
	}

	_, err = w.Write(byteResp)
	if err != nil {
		return err
	}

	return nil
}
