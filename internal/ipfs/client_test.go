package ipfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFile(t *testing.T) {
	path := "./test.test"
	name := "test.test"
	fu := CreateFile(path, name)
	assert.Equal(t, name, fu.Name)
	assert.Equal(t, path, fu.Path)
	assert.NotEqual(t, len(fu.CID), 0)
}
