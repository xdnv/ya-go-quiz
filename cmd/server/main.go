package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"internal/adapters/logger"
	"internal/app"

	"github.com/go-chi/chi/v5"
)

var sc app.ServerConfig

func handleGZIPRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(rw, r)
			return
		}

		logger.Info("srv-gzip: handling gzipped request")

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		defer gz.Close()
		body, err := io.ReadAll(gz)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		next.ServeHTTP(rw, r)
	})
}

func main() {
	//sync internal/logger upon exit
	defer logger.Sync()

	// create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// a WaitGroup for the goroutines to tell us they've stopped
	wg := sync.WaitGroup{}

	//Warning! do not run outside function, it will break tests due to flag.Parse()
	sc = app.InitServerConfig()

	// run `server` in it's own goroutine
	wg.Add(1)
	go server(ctx, &wg)

	// if err := run(); err != nil {
	// 	//logger.Error("Server error", zap.Error(err))
	// 	log.Fatal(err)
	// }

	// listen for ^C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	//fmt.Println("srv: received ^C - shutting down")
	logger.Info("srv: received ^C - shutting down")

	// tell the goroutines to stop
	//fmt.Println("srv: telling goroutines to stop")
	logger.Info("srv: telling goroutines to stop")
	cancel()

	// and wait for them to reply back
	wg.Wait()
	//fmt.Println("srv: shutdown")
	logger.Info("srv: shutdown")
}

func server(ctx context.Context, wg *sync.WaitGroup) {
	//execute to exit wait group
	defer wg.Done()

	logger.Info(fmt.Sprintf("srv: using endpoint %s", sc.Endpoint))
	logger.Info(fmt.Sprintf("srv: storage mode %v", sc.StorageMode))
	//logger.Info(fmt.Sprintf("srv: signed messaging=%v\n", signer.UseSignedMessaging()))

	mux := chi.NewRouter()
	//mux.Use(middleware.Logger)
	mux.Use(logger.LoggerMiddleware)
	mux.Use(handleGZIPRequests)
	//mux.Use(signer.HandleSignedRequests)
	//mux.Use(middleware.Compress(5, sc.CompressibleContentTypes...))

	mux.Get("/", index)
	mux.Get("/admin", adminPage)
	mux.Get("/quiz", quiz)
	mux.Post("/upload", uploadData)
	mux.Post("/command", handleCommand)
	mux.Post("/submit", submit)
	//mux.Get("/ping", pingDBServer)
	//mux.Post("/value/", requestMetricV2)
	//mux.Get("/value/{type}/{name}", requestMetricV1)
	//mux.Post("/update/", updateMetricV2)
	//mux.Post("/update/{type}/{name}/{value}", updateMetricV1)
	//mux.Post("/updates/", updateMetrics)

	// create a server
	srv := &http.Server{Addr: sc.Endpoint, Handler: mux}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			//fmt.Printf("Listen: %s\n", err)
			logger.Error(fmt.Sprintf("Listen: %s\n", err))
			//log.Fatal(err)
		}
	}()

	<-ctx.Done()
	//fmt.Println("srv: shutdown requested")
	logger.Info("srv: shutdown requested")

	// shut down gracefully with timeout of 5 seconds max
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ignore server error "Err shutting down server : context canceled"
	srv.Shutdown(shutdownCtx)

	// //save server state on shutdown
	// if sc.StorageMode == app.File {
	// 	err := stor.SaveState(sc.FileStoragePath)
	// 	if err != nil {
	// 		//fmt.Printf("srv: failed to save server state to [%s], error: %s\n", sc.FileStoragePath, err)
	// 		logger.Error(fmt.Sprintf("srv: failed to save server state to [%s], error: %s\n", sc.FileStoragePath, err))
	// 	}
	// }

	logger.Info("srv: server stopped")

}
