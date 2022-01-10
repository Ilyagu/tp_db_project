package middlware

import (
	"log"

	"github.com/valyala/fasthttp"
)

func ReponseMiddlwareAndLogger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		log.Printf("METHOD %s REMOTEADDR %s URL %s", ctx.Method(), ctx.RemoteAddr(), ctx.RequestURI())
		next(ctx)
	}
}
