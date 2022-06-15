package agent

import (
	"context"
	"log"
)

func Setup() error {
	log.Print("setting up")
	return nil
}

func Start(_ context.Context, _ chan error) error {
	log.Print("starting")
	return nil
}
