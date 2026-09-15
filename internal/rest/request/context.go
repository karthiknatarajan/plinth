package request

import (
	"context"
)

type key int

const (
	requestIDKey key = iota
)

// WithRequestID returns a copy of parent in which the request id value is set.
func WithRequestID(parent context.Context, v string) context.Context {
	return context.WithValue(parent, requestIDKey, v)
}
