package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/filipegorges/ports/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	err := app.Run(ctx)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}
}
