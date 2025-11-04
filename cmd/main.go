package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/root-N-root/webipfs/internal/http"
	"github.com/root-N-root/webipfs/internal/store"
	"github.com/root-N-root/webipfs/types"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	store.InitStore()
	con := types.NewConnector()
	go http.Run(ctx, con)
	go store.Run(ctx, con)
	<-ctx.Done()
}
