package store

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
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
	err := save(store)
	assert.NoError(t, err)
	data, err := load()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.EqualValues(t, data, store)
}

func TestLoadNoError(t *testing.T) {
	storeFilePath = getFilePath()
	defer os.Remove(storeFilePath)
	err := InitStore()
	assert.NoError(t, err)
	_, err = load()
	assert.NoError(t, err)
}

func TestSave(t *testing.T) {
	storeFilePath = getFilePath()
	defer os.Remove(storeFilePath)
	store := types.NewStore()
	err := save(store)
	assert.NoError(t, err)
}

func TestInitNoError(t *testing.T) {
	storeFilePath = getFilePath()
	defer os.Remove(storeFilePath)
	err := InitStore()
	assert.NoError(t, err)
}

func TestUploadFile(t *testing.T) {
	storeFilePath = getFilePath()
	defer os.Remove(storeFilePath)
	err := InitStore()
	assert.NoError(t, err)
	con := types.NewConnector()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	go Run(ctx, con)
	testFile := types.NewFileUpdate(types.FuwName("test.test"), types.FuwCid("test"), types.FuwPath("./test"))
	con.FileUpStoreChan <- testFile

	time.Sleep(1 * time.Second)
	store, err := load()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(store.Files))

	if len(store.Files) > 0 {
		assert.Equal(t, testFile.Name, store.Files[0].Name)
		assert.Equal(t, testFile.Path, store.Files[0].Path)
		assert.Equal(t, testFile.CID, store.Files[0].CID)
	}
}

func getFilePath() string {
	return path.Join(os.TempDir(), fmt.Sprintf("%d.json", time.Now().UnixNano()))
}
