package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nathann0209/microgo/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)
	helloHandler := handlers.NewHello(logger)
	goodbyeHandler := handlers.NewGoodbye(logger)
	postHandler := handlers.NewPost(logger)
	getHandler := handlers.NewGet(logger)
	fileHandler := handlers.NewFileHandler(logger)

	router := http.NewServeMux()
	router.Handle("/", helloHandler)
	router.Handle("/goodbye", goodbyeHandler)
	router.Handle("/post", postHandler)
	router.Handle("/get", getHandler)
	router.Handle("/pointers", fileHandler)
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	// Handle ListenandServe in go func "so with won't block"
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
		// http.ListenAndServe(":9090", sm)
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	logger.Println("Received terminate, graceful shutdown", sig)

	// time context
	// 30 seconds to attempt to gracefully shutdown but forcefully close down aftre 30 seconds
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
