package main

import (
	"context"
	logkit "github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func NewService(db *mongo.Client, logger logkit.Logger) Service {
	return &service{
		db: db,
		logger: logger,
	}
}

type service struct {
	db *mongo.Client
	logger logkit.Logger
}

type Service interface {
	GetHealthCheck() (string, error)
	GetProduct(Id string) (Product, error)
	StoreProduct(product Product) (Product, error)
}

type Product struct {
	Id 			string `json:"id"`
	Name 		string `json:"name"`
	Description string `json:"description"`
}

func (s *service) GetHealthCheck() (string, error) {
	return "ok", nil
}

func (s *service) GetProduct(Id string) (Product, error) {
	return Product{
		Id: Id,
		Name: "Test Product Name",
		Description: "Test Product Description",
	}, nil
}

func (s *service) StoreProduct(product Product) (Product, error) {
	NewProduct := struct {
		Name 		string `json:"name"`
		Description string `json:"description"`
	}{}

	NewProduct.Name = product.Name
	NewProduct.Description = product.Description

	collection := s.db.Database("products").Collection("catalogs")
	insertResult, err := collection.InsertOne(context.TODO(), NewProduct)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	InsertedID := insertResult.InsertedID.(primitive.ObjectID).Hex()

	s.logger.Log("Insert ID", InsertedID)

	return Product{
		Id: InsertedID,
		Name: product.Name,
		Description: product.Description,
	}, nil
}