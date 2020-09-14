package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{next: next, logger: logger}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) GetReviewList(ProductId string) (s []ProductReview, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetReviewList", "product_id", ProductId, "time", time.Since(begin), "err", err)
	}(time.Now())

	return mw.next.GetReviewList(ProductId)
}