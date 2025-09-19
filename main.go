package main

import (
	"context"
	"os"
	"os/signal"

	cli "github.com/grahamplata/roku-remote/cli/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cli.Run(ctx)
}
