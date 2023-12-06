package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	rabbit, err := newRabbit()
	if err != nil {
		log.Fatal(err)
	}

	mongo, err := newMongo()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	srvErrs := make(chan error, 1)
	go func() {
		srvErrs <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	shutdown := gracefulShutdown(srv, rabbit, mongo)

	select {
	case err := <-srvErrs:
		shutdown(err)
	case sig := <-quit:
		shutdown(sig)
	case err := <-rabbit.Err():
		shutdown(err)
	case err := <-mongo.Err():
		shutdown(err)
	}

	log.Println("Server exiting")
}

func gracefulShutdown(srv *http.Server, rabbit *Rabbit, mongo *Mongo) func(reason interface{}) {
	return func(reason interface{}) {
		log.Println("Server Shutdown:", reason)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Error Gracefully Shutting Down API:", err)
		}

		if err := rabbit.Disconnect(5 * time.Second); err != nil {
			log.Println("Error Gracefully Shutting Down Rabbit:", err)
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongo.Disconnect(ctx); err != nil {
			log.Println("Error Gracefully Shutting Down Mongo:", err)
		}
	}
}
