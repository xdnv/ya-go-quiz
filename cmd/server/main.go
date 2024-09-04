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
	"internal/ports/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var sc app.ServerConfig
var stor *storage.UniStorage

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

	stor = storage.NewUniStorage(&sc)
	defer stor.Close()

	//post-init unistorage actions
	err := stor.Bootstrap()
	if err != nil {
		logger.Fatal(fmt.Sprintf("srv: post-init bootstrap failed, error: %s\n", err))
	}

	// run `server` in it's own goroutine
	wg.Add(1)
	go server(ctx, &wg)

	// listen for ^C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logger.Info("srv: received ^C - shutting down")

	// tell the goroutines to stop
	logger.Info("srv: telling goroutines to stop")
	cancel()

	// and wait for them to reply back
	wg.Wait()
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
	mux.Use(middleware.Compress(5, sc.CompressibleContentTypes...))

	mux.Get("/", index)
	mux.Get("/admin", adminPage)
	mux.Get("/quiz/{id}", quiz)
	mux.Get("/login", authPage)
	mux.Get("/logout", logout)
	mux.Get("/results/{id}", handleResults)
	mux.Post("/login", auth)
	mux.Post("/upload", uploadData)
	mux.Post("/command", handleCommand)
	mux.Post("/submit", submit)

	// create a server
	srv := &http.Server{Addr: sc.Endpoint, Handler: mux}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(fmt.Sprintf("Listen: %s\n", err))
			//log.Fatal(err)
		}
	}()

	<-ctx.Done()
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
