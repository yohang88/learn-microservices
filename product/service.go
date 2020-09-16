package main

import (
	"context"
	logkit "github.com/go-kit/kit/log"
	reviewpb "github.com/yohang88/learn-microservices/review/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func NewService(db *mongo.Client, logger logkit.Logger, review reviewpb.ReviewServiceClient) Service {
	return &service{
		db: db,
		logger: logger,
		review: review,
	}
}

type service struct {
	db *mongo.Client
	logger logkit.Logger
	review reviewpb.ReviewServiceClient
}

type Service interface {
	GetHealthCheck() (string, error)
	GetProductList() ([]Product, error)
	GetProduct(Id string) (Product, error)
	StoreProduct(product Product) (Product, error)
}

type Product struct {
	Id 			string   `json:"id" bson:"_id"`
	Name 		string   `json:"name"`
	Description string 	 `json:"description"`
	Reviews 	[]Review `json:"reviews,omitempty"`
}

type Review struct {
	Id 			string `json:"id"`
	Content 	string `json:"content"`
}


func (s *service) GetHealthCheck() (string, error) {
	return "ok", nil
}

func (s *service) GetProductList() ([]Product, error) {
	var results []Product

	collection := s.db.Database("products").Collection("catalogs")

	findOptions := options.Find()

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)

	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem Product
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return results, nil
}

func (s *service) GetProduct(Id string) (Product, error) {
	docID, err := primitive.ObjectIDFromHex(Id)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	filter := bson.M{"_id": docID}

	var product Product

	collection := s.db.Database("products").Collection("catalogs")

	err = collection.FindOne(context.TODO(), filter).Decode(&product)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	reviewResponse, errResponse := s.review.GetReviewList(context.Background(), &reviewpb.GetReviewListRequest{
		ProductId: Id,
	})

	if errResponse != nil {
		s.logger.Log("message", "Could not connect to Review Service")

		return Product{
			Id: Id,
			Name: product.Name,
			Description: product.Description,
		}, nil
	}

	var reviews []Review

	for _, record := range reviewResponse.Reviews {
		reviews = append(reviews, Review{
			Id: record.GetId(),
			Content: record.GetContent(),
		})
	}

	return Product{
		Id: Id,
		Name: product.Name,
		Description: product.Description,
		Reviews: reviews,
	}, nil
}

func (s *service) StoreProduct(product Product) (Product, error) {
	NewProduct := struct{
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