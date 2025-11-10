package main

import (
	"context"
	"fmt"
	"log"
	"microservice/app"
	"microservice/internal/server/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Service struct {
	service *app.App
	http    http.IHttpServer
}

func main() {
	service := New()
	service.start()
}

func New() *Service {
	return &Service{}
}

func (a *Service) start() {
	fmt.Printf("\n[service] starting...\n")

	a.service = app.New()
	a.service.Init()

	if a.service.Config().Debug == true {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	a.http = http.New(
		a.service.Registry(),
		a.service.Locale(),
		a.service.Repositories(),
		a.service.HttpHandlers(),
	)

	a.http.Init()
	a.http.SetRoutes()
	a.http.Start()

	// NOTE: init the gRPC server here as needed

	fmt.Printf("[service] started\n")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-ch

	a.stop()
}

// stop is prioritized to act as minimal graceful shutdown
func (a *Service) stop() {
	fmt.Printf("\n[usecase] shutting down...\n")
	defer fmt.Printf("\n[usecase] shutdown successfully\n")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.service.Config().StopTimeout)*time.Second)
	defer cancel()

	a.http.Stop(ctx)
	a.service.DB().Stop()
	a.service.Logger().Stop()
}
