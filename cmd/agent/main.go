package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adampointer/agent/internal/app/agent"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	errC := make(chan error, 1)

	if err := agent.Setup(); err != nil {
		log.Fatal(err)
	}

	if err := agent.Start(ctx, errC); err != nil {
		log.Fatal(err)
	}

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	case <-terminate:
		log.Print("shutting down")
		cancel()
	case err := <-errC:
		log.Fatal(err)
	}
}
