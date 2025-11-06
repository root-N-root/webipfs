package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/root-N-root/webipfs/internal/http"
	"github.com/root-N-root/webipfs/internal/ipfs"
	"github.com/root-N-root/webipfs/internal/store"
	"github.com/root-N-root/webipfs/types"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	store.InitStore()
	con := types.NewConnector()
	client, err := ipfs.Initialize(ctx, con)
	if err != nil {
		log.Fatal(err)
	}
	go http.Run(ctx, con, client)
	go store.Run(ctx, con)
	<-ctx.Done()
}
