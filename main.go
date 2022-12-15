package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func createStop() (chan os.Signal, func()) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	return stop, func() {
		close(stop)
	}
}

func startServer(s *http.Server) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	log.Print("listening...")
	if err := s.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
		log.Print("http server closed")
	}
}

func shutdownServer(ctx context.Context, s *http.Server) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Print("server shutdown")
}

func main() {
	stop, closeCh := createStop()
	defer closeCh()

	srv := &http.Server{
		Addr:    ":9000",
		Handler: nil,
	}

	go startServer(srv)

	<-stop

	shutdownServer(context.Background(), srv)
}
