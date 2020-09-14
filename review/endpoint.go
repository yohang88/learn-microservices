package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetReviewListEndpoint endpoint.Endpoint
}

func makeServerEndpoints(svc Service) Endpoints {
	return Endpoints{
		GetReviewListEndpoint: makeGetReviewListEndpoint(svc),
	}
}

func makeGetReviewListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReviewListRequest)

		result, _ := svc.GetReviewList(req.ProductId)

		return GetReviewListResponse(result), nil
	}
}