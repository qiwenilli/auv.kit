package utils

import (
	"context"
	"fmt"
	"net/http"

	auvhttp "github.com/qiwenilli/auv.kit/internal/http"
)

func TraceIdForContext(ctx context.Context) string {
	return fmt.Sprintf("%s", ctx.Value(auvhttp.TraceIdName))
}

func TraceId(r *http.Request) string {
	return TraceIdForContext(r.Context())
}