package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := makeServerEndpoints(s)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/health").Handler(httptransport.NewServer(
		e.GetHealthCheckEndpoint,
		decodeGetHealthCheckRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/products").Handler(httptransport.NewServer(
		e.GetProductListEndpoint,
		decodeGetProductListRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/products/{id}").Handler(httptransport.NewServer(
		e.GetProductEndpoint,
		decodeGetProductRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/products").Handler(httptransport.NewServer(
		e.StoreProductEndpoint,
		decodeStoreProductRequest,
		encodeResponse,
		options...,
	))

	return r
}

type GetHealthCheckRequest struct {}

type GetHealthCheckResponse struct {
	Status string `json:"status"`
}

type GetProductListRequest struct {}

type GetProductListResponse []Product

type GetProductRequest struct {
	Id string `json:"id"`
}

type GetProductResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type StoreProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type StoreProductResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func decodeGetHealthCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetHealthCheckRequest
	return req, nil
}

func decodeGetProductListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetProductListRequest
	return req, nil
}

func decodeGetProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	var req GetProductRequest

	req.Id = vars["id"]

	return req, nil
}

func decodeStoreProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req StoreProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
