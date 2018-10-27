package middleware

import (
	"context"
	"net/http"

	"github.com/sahandhnj/apiclient/util"
)

type key int

const requestIDKey key = 0

func newContextWithRequestID(ctx context.Context, req *http.Request) context.Context {
	reqID := req.Header.Get("X-Request-ID")
	if reqID == "" {
		reqID = util.UUID()
	}

	return context.WithValue(ctx, requestIDKey, reqID)
}

func requestIDFromContext(ctx context.Context) string {
	if ctx.Value(requestIDKey) != nil {
		return ctx.Value(requestIDKey).(string)
	}

	return ""
}
