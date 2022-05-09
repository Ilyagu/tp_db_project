package middlware

import (
	"log"
	"net/http"
)

func ReponseMiddlwareAndLogger() (mw func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			log.Printf("METHOD %s REMOTEADDR %s URL %s", r.Method, r.RemoteAddr, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	}
}
