package types

import (
	"context"
	"golang.org/x/time/rate"
	"io"
)

type RateLimitWriter struct {
	writer  io.Writer
	limit   int
	limiter *rate.Limiter
	context context.Context
	sktInfo *SocketInfo
}
