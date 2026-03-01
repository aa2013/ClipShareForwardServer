package ratelimiter

import (
	"context"
	"golang.org/x/time/rate"
	"io"
)

type RateLimitWriter struct {
	writer          io.Writer
	limit           int
	limiter         *rate.Limiter
	context         context.Context
	afterWriteBytes func(cnt int)
}

func NewLimitedWriterCallBack(writer io.Writer, limit int, afterWriteBytes func(cnt int)) *RateLimitWriter {
	limiter := rate.NewLimiter(rate.Limit(limit), limit)
	return &RateLimitWriter{
		writer:          writer,
		limiter:         limiter,
		limit:           limit,
		context:         context.Background(),
		afterWriteBytes: afterWriteBytes,
	}
}
func NewLimitedWriter(writer io.Writer, limit int) *RateLimitWriter {
	return NewLimitedWriterCallBack(writer, limit, nil)
}
func (w *RateLimitWriter) Write(p []byte) (int, error) {
	n := 0
	pLen := len(p)
	for n < pLen {
		// 计算最大速率
		chunkSize := w.limiter.Burst()
		if chunkSize <= 0 {
			chunkSize = pLen - n
		} else {
			if chunkSize > pLen-n {
				chunkSize = pLen - n
			}
			// 等待
			err := w.limiter.WaitN(w.context, chunkSize)
			if err != nil {
				return n, err
			}
		}
		// 写入数据
		cnt, err := w.writer.Write(p[n : n+chunkSize])
		if w.afterWriteBytes != nil {
			w.afterWriteBytes(cnt)
		}
		if err != nil {
			return n, err
		}
		if cnt < chunkSize {
			return n, nil
		}
		n += cnt
	}

	return n, nil
}
func (w *RateLimitWriter) UpdateLimit(limit int) {
	w.limit = limit
	w.limiter = rate.NewLimiter(rate.Limit(limit), limit)
}
