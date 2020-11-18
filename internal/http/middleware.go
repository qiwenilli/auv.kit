package http

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/juju/ratelimit"
	log "github.com/sirupsen/logrus"

	"github.com/qiwenilli/auv.kit/utils"
)

func MiddlewareResponseTime(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inTime := time.Now()
		h.ServeHTTP(w, r)
		log.Infof("%s %s %v %s %s", r.Method, r.URL, r.Header, utils.RemoteIP(r), time.Now().Sub(inTime))
	})
}

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

func MiddlewareTraceId(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var traceId string
		if headerTraceIds, ok := r.Header[utils.TraceIdName]; !ok {
			traceId = fmt.Sprintf("%d", uuid.New().ID())
			r.Header.Add(utils.TraceIdName, traceId)
		} else {
			traceId = strings.Join(headerTraceIds, " ")
		}
		//
		ctx := context.WithValue(r.Context(), utils.TraceIdName, traceId)
		r = r.WithContext(ctx)

		w.Header().Add(utils.TraceIdName, traceId)

		h.ServeHTTP(w, r)
	})
}

var (
	limiterOnce   sync.Once
	limiterEntity *limiter
)

type limiter struct {
	everySecondQps int64
	waitTimeout    time.Duration
	ipWhiteList    []string
	bucker         *ratelimit.Bucket
}

func NewMiddlewareRlimit(everySecondQps int64) *limiter {
	limiterOnce.Do(func() {
		limiterEntity = &limiter{
			everySecondQps: everySecondQps,
			bucker:         ratelimit.NewBucket(time.Second, everySecondQps),
		}
	})
	return limiterEntity
}

func (l *limiter) SetIpWhiteList(ips []string) {
	for _, ip := range ips {
		ip = strings.Replace(ip, "*", `[0-9]*`, -1)
		ip = strings.Replace(ip, ".", `\.`, -1)
		ip = `^` + ip
		l.ipWhiteList = append(l.ipWhiteList, ip)
	}
	log.Debug(l.ipWhiteList)
}

func (l *limiter) SetWaitTimeout(t time.Duration) {
	l.waitTimeout = t
}

func (l *limiter) inWhiteList(visitorIp string) bool {
	if len(l.ipWhiteList) > 0 {
		for _, ip := range l.ipWhiteList {
			if ok, _ := regexp.MatchString(ip, visitorIp); ok {
				return true
			}
		}
		return false
	}
	return true
}

func (l *limiter) MiddlewareRlimit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		visitorIp := utils.RemoteIP(r)
		if l.inWhiteList(visitorIp) {
			log.Infof("visitorIp=%s in white list", visitorIp)
			h.ServeHTTP(w, r)
			return
		}
		if ok := l.bucker.WaitMaxDuration(1, l.waitTimeout); ok {
			h.ServeHTTP(w, r)
		} else {
			log.Warnf("TooManyRequests: code=%d ratelimit=%d/second", http.StatusTooManyRequests, l.everySecondQps)
			w.WriteHeader(http.StatusTooManyRequests)
		}
	})
}
