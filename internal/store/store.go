package store

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/root-N-root/webipfs/types"
)

var storeFilePath = types.STORE_FILE_PATH

func Run(ctx context.Context, con *types.Connector) {
	for {
		select {
		case fu := <-con.FileUpStoreChan:
			err := updateInStore(fu)
			if err != nil {
				log.Println("store.run:", err)
			}
		case <-ctx.Done():
			break
		}
	}
}

func updateInStore(fu types.FileUpdate) error {
	store, err := load()
	if err != nil {
		return err
	}
	store.UpdateFile(fu)
	return save(store)

}

func InitStore() error {
	var err error
	needFill := false
	if _, err = os.Stat(storeFilePath); err != nil {
		needFill = true
		_, err = os.Create(storeFilePath)
		if err != nil {
			return err
		}
	}
	if !needFill {
		if _, err = load(); err != nil {
			needFill = true
		}
	}
	if needFill {
		return save(types.NewStore())
	}
	return err
}

func load() (types.Store, error) {
	var store types.Store
	var err error
	if _, err = os.Stat(storeFilePath); err != nil {
		return store, err
	}
	file, err := os.Open(storeFilePath)
	if err != nil {
		return store, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&store)
	return store, err
}

func save(store types.Store) error {
	data, err := json.MarshalIndent(store, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(storeFilePath, data, 0644)
}
