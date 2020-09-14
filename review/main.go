package main

import (
	"context"
	"fmt"
	logkit "github.com/go-kit/kit/log"
	reviewpb "github.com/yohang88/learn-microservices/review/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var logger logkit.Logger
	{
		logger = logkit.NewLogfmtLogger(os.Stderr)
		logger = logkit.With(logger, "timestamp", logkit.DefaultTimestampUTC)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	var s Service
	{
		s = NewService(client, logger)
		s = LoggingMiddleware(logger)(s)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:50000")

		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		gRPCServer := grpc.NewServer()

		endpoints := MakeGRPCServer(ctx, makeServerEndpoints(s))

		reviewpb.RegisterReviewServiceServer(gRPCServer, endpoints)

		log.Println("Review Service is listening on port 50000...")
		errs <- gRPCServer.Serve(listener)
	}()

	logger.Log("exit", <-errs)
}

