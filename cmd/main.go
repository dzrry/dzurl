package main

import (
	"fmt"
	rediss "github.com/dzrry/dzurl/repo/redis"
	"github.com/dzrry/dzurl/service"
	"github.com/dzrry/dzurl/transport"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	rr, err := rediss.NewRepo("localhost", "6379", "")
	if err != nil {
		log.Fatal(err)
	}
	srvc := service.NewRedirectService(rr)
	handler := transport.NewHandler(srvc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{key}", handler.LoadRedirect)
	r.Post("/", handler.StoreRedirect)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8080")
		errs <- http.ListenAndServe("localhost:8080", r)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}
