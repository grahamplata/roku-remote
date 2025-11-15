package main

import (
	"context"
	"os"
	"os/signal"

	cli "github.com/grahamplata/roku-remote/cli/cmd"
)

// version is set via ldflags during build
var version = "dev"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cli.Run(ctx, version)
}
