package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	READ_TIMEOUT           = 10
	WRITE_TIMEOUT          = 10
	IDLE_TIMEOUT           = 30
	LOG_FILE_NAME          = "server_log.log"
	DEFAULT_SERVER_ADDRESS = "8888"
)

var (
	listenAddr string
	healthy    int32
)

func main() {
	flag.StringVar(&listenAddr, "address", ":"+DEFAULT_SERVER_ADDRESS, "http service address")
	flag.Parse()

	logFile, err := os.OpenFile(LOG_FILE_NAME, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	defer logFile.Close()

	logger := log.New(logFile, "http: ", log.LstdFlags)

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      corsMiddlware()(logMiddleware(logger)(routeHandler())),
		ErrorLog:     logger,
		ReadTimeout:  READ_TIMEOUT * time.Second,
		WriteTimeout: WRITE_TIMEOUT * time.Second,
		IdleTimeout:  IDLE_TIMEOUT * time.Second,
	}

	doneChannel := make(chan bool)
	quitChannel := make(chan os.Signal, 1)

	signal.Notify(quitChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quitChannel
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)

		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}

		close(doneChannel)
	}()

	logger.Println("Server is ready to handle requests at ", listenAddr)

	// Atomically update our health state indicator to 'healthy'
	atomic.StoreInt32(&healthy, 1)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	// If we receive a signal via the done channel, we log the event:
	<-doneChannel
	logger.Println("Server stopped")
}

func corsMiddlware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add CORS headers to the response
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			// Handle preflight requests (OPTIONS)
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
