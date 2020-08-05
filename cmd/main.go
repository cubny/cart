package main

import (
	"context"
	"flag"
	"github.com/cubny/cart/internal/auth"
	"github.com/cubny/cart/internal/handler"
	"github.com/cubny/cart/internal/service"
	"github.com/cubny/cart/internal/storage/sqlite3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var (
		optsAddr = flag.String("addr", ":5120", "HTTP bind address")
		dataPath = flag.String("data", "./data/db.sqlite3", "Path to the DB file")
	)
	flag.Parse()

	file, err := os.Open(*dataPath)
	if err != nil {
		log.Fatalf("cannot open dbPath file: %s", err)
	}

	storage, err := sqlite3.New(file)
	if err != nil {
		log.Fatalf("sqlite3: %s", err)
	}

	service, err := service.New(storage)
	if err != nil {
		log.Fatalf("cannot create service, %s", err)
	}

	authClient := auth.New()
	handler, err := handler.New(service, authClient)
	if err != nil {
		log.Fatalf("cannot create handler, %s", err)
	}

	srv := http.Server{
		Addr:    *optsAddr,
		Handler: handler,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Printf("shuting down the http server...")
		idleCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := srv.Shutdown(idleCtx); err != nil {
			panic(err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("HTTP Server starting %s", *optsAddr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
	<-idleConnsClosed
}
