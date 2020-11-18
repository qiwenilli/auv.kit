package utils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const (
	TraceIdName = "X-Auv-TraceId"
)

func traceIdForContext(ctx context.Context) string {
	return fmt.Sprintf("%s", ctx.Value(TraceIdName))
}

func TraceId(r *http.Request) string {
	return traceIdForContext(r.Context())
}

func RemoteIP(r *http.Request) string {
	remoteAddr := strings.TrimSpace(r.RemoteAddr)
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return ""
	}
	return host
}
