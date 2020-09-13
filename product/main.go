package main

import (
	"context"
	"flag"
	"fmt"
	logkit "github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

