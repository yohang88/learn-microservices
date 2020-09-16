package main

import (
	"context"
	"flag"
	"fmt"
	logkit "github.com/go-kit/kit/log"
	reviewpb "github.com/yohang88/learn-microservices/review/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)
	flag.Parse()

	var logger logkit.Logger
	{
		logger = logkit.NewLogfmtLogger(os.Stderr)
		logger = logkit.With(logger, "timestamp", logkit.DefaultTimestampUTC)
	}

	mongoConnection := os.Getenv("MONGO_CONNECTION")

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoConnection))
	if err != nil {
		log.Fatalf("Cannot connect to MongoDB: %v", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoClient.Connect(ctx)

	var review reviewpb.ReviewServiceClient
	{
		grpcConn, err := grpc.Dial("localhost:50000", grpc.WithInsecure())
		client := reviewpb.NewReviewServiceClient(grpcConn)

		if err != nil {
			log.Fatalf("Could not connect: %v", err)
		}

		defer grpcConn.Close()

		review = client
	}


	var s Service
	{
		s = NewService(mongoClient, logger, review)
		s = LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = MakeHTTPHandler(s, logkit.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}

