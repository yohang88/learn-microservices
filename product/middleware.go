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

func (mw loggingMiddleware) GetHealthCheck() (s string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetHealthCheck", "time", time.Since(begin), "err", err)
	}(time.Now())

	return mw.next.GetHealthCheck()
}

func (mw loggingMiddleware) GetProduct(id string) (product Product, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetProduct", "id", id, "time", time.Since(begin), "err", err)
	}(time.Now())

	return mw.next.GetProduct(id)
}

func (mw loggingMiddleware) StoreProduct(req Product) (res Product, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "StoreProduct", "time", time.Since(begin), "err", err)
	}(time.Now())

	return mw.next.StoreProduct(req)
}