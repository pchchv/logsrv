package logging

import "net/http"

type LogMiddleware struct {
	Next http.Handler
}

func NewLogMiddleware(next http.Handler) *LogMiddleware {
	return &LogMiddleware{
		Next: next,
	}
}
