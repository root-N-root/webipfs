package main

import (
	"github.com/root-N-root/webipfs/internal/http"
	"github.com/root-N-root/webipfs/internal/store"
	"github.com/root-N-root/webipfs/types"
)

func main() {
	store.InitStore()
	fuCh, msgCh := types.NewChans()
	// go store.Run(fuch, msgc)
	go http.Run(fuCh, msgCh)
}
