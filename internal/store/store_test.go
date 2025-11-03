package store

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/root-N-root/webipfs/types"
)

func TestFull(t *testing.T) {
	storeFilePath = getFilePath()
	defer os.Remove(storeFilePath)
	store := types.NewStore()
	store.AddFile(types.NewFileUpdate(types.FuwCid("test")))
	err := Save(store)
	assert.Nil(t, err)
	data, err := Load()
	assert.Nil(t, err)
	assert.NotEmpty(t, data)
	assert.EqualValues(t, data, store)
}

func TestLoadNoError(t *testing.T) {
	storeFilePath = getFilePath()
	defer os.Remove(storeFilePath)
	err := InitStore()
	assert.Nil(t, err)
	_, err = Load()
	assert.Nil(t, err)
}

func TestSave(t *testing.T) {
	storeFilePath = getFilePath()
	defer os.Remove(storeFilePath)
	store := types.NewStore()
	err := Save(store)
	assert.Nil(t, err)
}

func TestInitNoError(t *testing.T) {
	storeFilePath = getFilePath()
	defer os.Remove(storeFilePath)
	err := InitStore()
	assert.Nil(t, err)
}

func getFilePath() string {
	return path.Join(os.TempDir(), fmt.Sprintf("%d.json", time.Now().UnixNano()))
}
