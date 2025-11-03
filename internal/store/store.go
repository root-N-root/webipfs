package store

import (
	"encoding/json"
	"os"

	"github.com/root-N-root/webipfs/types"
)

var storeFilePath = types.STORE_FILE_PATH

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
		if _, err = Load(); err != nil {
			needFill = true
		}
	}
	if needFill {
		return Save(types.NewStore())
	}
	return err
}

func Load() (types.Store, error) {
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

func Save(store types.Store) error {
	data, err := json.MarshalIndent(store, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(storeFilePath, data, 0644)
}
