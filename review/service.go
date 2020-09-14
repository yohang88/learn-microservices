package main

import (
	"context"
	logkit "github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	GetReviewList(ProductId string) ([]ProductReview, error)
}

type ProductReview struct {
	Id 			string `json:"id" bson:"_id"`
	ProductId   string `json:"product_id" bson:"_id"`
	Content     string `json:"content"`
}

func (s *service) GetReviewList(ProductId string) ([]ProductReview, error) {
	docProductID, err := primitive.ObjectIDFromHex(ProductId)

	var results []ProductReview

	filter := bson.M{"product_id": docProductID}

	collection := s.db.Database("product_reviews").Collection("reviews")

	findOptions := options.Find()

	cur, err := collection.Find(context.TODO(), filter, findOptions)

	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem ProductReview
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