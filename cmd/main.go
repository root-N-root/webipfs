package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	// Import IPFS plugins to register datastore types
	_ "github.com/ipfs/go-ds-flatfs"
	_ "github.com/ipfs/go-ds-leveldb"
	_ "github.com/ipfs/go-ds-measure"
	
	"github.com/ipfs/kubo/plugin/loader"

	"github.com/root-N-root/webipfs/internal/http"
	"github.com/root-N-root/webipfs/internal/ipfs"
	"github.com/root-N-root/webipfs/internal/store"
	"github.com/root-N-root/webipfs/types"
)

func init() {
	// Try to load IPFS plugins at initialization time
	plugins, err := loader.NewPluginLoader("")
	if err != nil {
		log.Printf("Failed to load IPFS plugin loader: %v", err)
	} else {
		if err := plugins.Initialize(); err != nil {
			log.Printf("Failed to initialize IPFS plugins: %v", err)
		} else {
			if err := plugins.Inject(); err != nil {
				log.Printf("Failed to inject IPFS plugins: %v", err)
			}
		}
	}
}

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
