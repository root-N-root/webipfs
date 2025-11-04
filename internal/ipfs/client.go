package ipfs

import (
	"github.com/root-N-root/webipfs/types"
)

//TODO::init dht and other

func CreateFile(filepath string, filename string) types.FileUpdate {
	//TODO:: create fu
	return types.NewFileUpdate(types.FuwName(filename), types.FuwPath(filepath), types.FuwCid("test"))
}
