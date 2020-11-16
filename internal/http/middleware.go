package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	// log "github.com/sirupsen/logrus"
)

func MiddlewareForCrossDomain(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "X-Auv-Cors")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		//
		h.ServeHTTP(w, r)
	})
}

const TraceIdName = "X-Auv-TraceId"

func MiddlewareTraceId(h http.Handler) http.Handler {

	var headerRWLock = sync.RWMutex{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var traceId string
		headerRWLock.Lock()
		if headerTraceIds, ok := r.Header[TraceIdName]; !ok {
			traceId = fmt.Sprintf("%d", uuid.New().ID())
			r.Header.Add(TraceIdName, traceId)
		} else {
			traceId = strings.Join(headerTraceIds, " ")
		}
		headerRWLock.Unlock()
		//
		ctx := context.WithValue(r.Context(), TraceIdName, traceId)
		r = r.WithContext(ctx)

		w.Header().Add(TraceIdName, traceId)

		h.ServeHTTP(w, r)
	})
}
