package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetHealthCheckEndpoint endpoint.Endpoint
	GetProductEndpoint endpoint.Endpoint
}

func makeServerEndpoints(svc Service) Endpoints {
	return Endpoints{
		GetHealthCheckEndpoint: makeHealthCheckEndpoint(svc),
		GetProductEndpoint: makeGetProductEndpoint(svc),
	}
}

func makeHealthCheckEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		status, _ := svc.GetHealthCheck()

		return GetHealthCheckResponse{status}, nil
	}
}

func makeGetProductEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetProductRequest)

		result, _ := svc.GetProduct(req.Id)

		return GetProductResponse{
			Id: result.Id,
			Name: result.Name,
			Description: result.Description,
		}, nil
	}
}