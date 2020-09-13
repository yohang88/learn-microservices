package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetHealthCheckEndpoint endpoint.Endpoint
	GetProductEndpoint endpoint.Endpoint
	StoreProductEndpoint endpoint.Endpoint
}

func makeServerEndpoints(svc Service) Endpoints {
	return Endpoints{
		GetHealthCheckEndpoint: makeHealthCheckEndpoint(svc),
		GetProductEndpoint: makeGetProductEndpoint(svc),
		StoreProductEndpoint: makeStoreProductEndpoint(svc),
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

func makeStoreProductEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StoreProductRequest)

		result, _ := svc.StoreProduct(Product{
			Name: req.Name,
			Description: req.Description,
		})

		return StoreProductResponse{
			Id: result.Id,
			Name: result.Name,
			Description: result.Description,
		}, nil
	}
}