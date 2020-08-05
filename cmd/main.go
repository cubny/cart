package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cubny/cart/internal/auth"
	"github.com/cubny/cart/internal/handler"
	"github.com/cubny/cart/internal/service"
	"github.com/cubny/cart/internal/storage/sqlite3"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		optsAddr    = flag.String("addr", ":8080", "HTTP bind address")
		metricsAddr = flag.String("metricsAddr", ":8081", "Metrics HTTP bind address")
		dataPath    = flag.String("data", "/app/data/cart.db", "Path to the sqlite3 data file")
		migrate     = flag.Bool("migrate", false, "if migrate is set the migration will be performed")
	)
	flag.Parse()

	// Check if data file exist
	_, err := os.Stat(*dataPath)
	if os.IsNotExist(err) && !*migrate {
		log.Fatalf("file %s does not exist. To create a new database run the process with -migrate flag", *dataPath)
	}

	// open or create the data file
	dbFile, err := os.OpenFile(*dataPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("cannot open data file %s: %s", *dataPath, err)
	}

	storage, err := sqlite3.New(dbFile)
	if err != nil {
		log.Fatalf("sqlite3: %s", err)
	}

	if *migrate {
		log.Println("migrating started")
		if err := storage.Migrate(); err != nil {
			log.Fatalf("migration failed, %s", err)
		}
		log.Println("migrating finished.")
		log.Println("run the program again without the migrate flag to start the server")
		os.Exit(0)
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

	go func() {
		log.Debugf("starting metrics server %s", *metricsAddr)
		log.Fatal(http.ListenAndServe(*metricsAddr, promhttp.Handler()))
	}()

	log.Printf("HTTP Server starting %s", *optsAddr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
	<-idleConnsClosed
}
