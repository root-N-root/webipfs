package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFileUpdate = FileUpdate{
	CID:      "test-cid",
	Name:     "test-name",
	Peers:    1,
	Progress: 0.0,
	Type:     "test-type",
	Status:   StatusComplete,
}

func TestStoreUpdate(t *testing.T) {
	s := NewStore()
	f := NewFileUpdate(FuwCid(testFileUpdate.CID))
	s.Files = append(s.Files, f)
	f_updated := NewFileUpdate(FuwCid(testFileUpdate.CID), FuwName(testFileUpdate.Name))
	s.UpdateFile(f_updated)
	assert.Equal(t, len(s.Files), 1)
	assert.Equal(t, s.Files[0].Name, testFileUpdate.Name)
}

func TestFileUpdateWithFuncs(t *testing.T) {

	fu := NewFileUpdate(
		FuwCid(testFileUpdate.CID),
		FuwName(testFileUpdate.Name),
		FuwPeers(testFileUpdate.Peers),
		FuwProgress(testFileUpdate.Progress),
		FuwType(testFileUpdate.Type),
		FuwStatus(testFileUpdate.Status),
	)
	assert.Equal(t, fu.CID, testFileUpdate.CID)
	assert.Equal(t, fu.Name, testFileUpdate.Name)
	assert.Equal(t, fu.Peers, testFileUpdate.Peers)
	assert.Equal(t, fu.Progress, testFileUpdate.Progress)
	assert.Equal(t, fu.Type, testFileUpdate.Type)
	assert.Equal(t, fu.Status, testFileUpdate.Status)
}
