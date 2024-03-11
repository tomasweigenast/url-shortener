package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"tomasweigenast.com/url-shortener/database"
	"tomasweigenast.com/url-shortener/server/endpoints"
	"tomasweigenast.com/url-shortener/server/middleware"
	"tomasweigenast.com/url-shortener/taskmanager"
)

var mux *http.ServeMux
var server *http.Server

const defaultPort = 9999

func Run() {
	port := getPort()
	mux = http.NewServeMux()
	mux.Handle("PUT /register", middleware.Chain(endpoints.Register, middleware.RequestID, middleware.Logging))
	mux.Handle("POST /login", middleware.Chain(endpoints.Login, middleware.RequestID, middleware.Logging))
	mux.Handle("GET /account/", middleware.Chain(endpoints.GetAccount, middleware.RequestID, middleware.Logging, middleware.Auth))
	mux.Handle("GET /account/links", middleware.Chain(endpoints.GetLinks, middleware.RequestID, middleware.Logging, middleware.Auth))
	mux.Handle("PUT /account/links", middleware.Chain(endpoints.CreateLink, middleware.RequestID, middleware.Logging, middleware.Auth))
	mux.Handle("GET /account/links/{id}", middleware.Chain(endpoints.GetLink, middleware.RequestID, middleware.Logging, middleware.Auth))
	mux.Handle("DELETE /account/links/{id}", middleware.Chain(endpoints.DeleteLink, middleware.RequestID, middleware.Logging, middleware.Auth))
	mux.Handle("GET /{url}", middleware.Chain(endpoints.Url, middleware.RequestID, middleware.Logging))

	server = &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", port),
	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Starting server..")
	database.ConnectDatabase()
	taskmanager.Start()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("could not start server: %s", err)
			}
		}
	}()
	log.Printf("Server started at %s \n", server.Addr)

	<-done
	log.Println("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		taskmanager.Stop()
		database.CloseDatabase()
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("unable to shudown server gracefully: %s", err)
	}

	log.Println("Server stopped successfully.")
}

func getPort() int64 {
	if portStr, ok := os.LookupEnv("PORT"); ok {
		port, err := strconv.ParseInt(portStr, 10, 32)
		if err != nil {
			log.Fatalf("invalid port, must be an integer: %s", err)
		}

		return port
	}

	return defaultPort
}
