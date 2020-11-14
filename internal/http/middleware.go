package http

import (
	"net/http"
)

func MiddlewareForCrossDomain(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add(" Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "X-Auv-Cors")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		//
		h.ServeHTTP(w, r)
	})
}
