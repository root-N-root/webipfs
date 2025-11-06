package types

const PORT = 3000
const STORE_FILE_PATH = "./webipfs-store.json"
const FILES_DIR = "./files"

type FileService interface {
	AddFile(string, string) (FileUpdate, error)
}

type Connector struct {
	FileUpHttpChan  chan FileUpdate
	FileUpStoreChan chan FileUpdate
	MsgChan         chan string
}

func NewConnector() *Connector {
	return &Connector{
		FileUpHttpChan:  make(chan FileUpdate),
		FileUpStoreChan: make(chan FileUpdate),
		MsgChan:         make(chan string),
	}
}

func (c *Connector) SendFileUp(fu FileUpdate) {
	c.FileUpHttpChan <- fu
	c.FileUpStoreChan <- fu
}

type Store struct {
	Files []FileUpdate `json:"files"`
}

func NewStore() Store {
	return Store{}
}

func (store *Store) AddFile(file FileUpdate) {
	store.Files = append(store.Files, file)
}

func (store *Store) RemoveFile(cid string) *FileUpdate {
	for index, file := range store.Files {
		if file.CID == cid {
			store.Files = append(store.Files[:index], store.Files[index+1:]...)
			return &file
		}
	}
	return nil
}

func (store *Store) UpdateFile(file FileUpdate) {
	for i, fileInStore := range store.Files {
		if file.CID == fileInStore.CID {
			store.Files[i] = file
			return
		}
	}
	store.AddFile(file)
}

type FileStatus string

const (
	StatusQueued      FileStatus = "queued"
	StatusDownloading FileStatus = "downloading"
	StatusSeeding     FileStatus = "seeding"
	StatusComplete    FileStatus = "complete"
	StatusError       FileStatus = "error"
)

type FileUpdate struct {
	CID      string     `json:"cid"`
	Name     string     `json:"name,omitempty"`
	Peers    int        `json:"peers"`
	Progress float64    `json:"progress"`
	Type     string     `json:"type"`
	Status   FileStatus `json:"file_status"`
	Path     string     `json:"path"`
}

func NewFileUpdate(fuWiths ...FUWith) FileUpdate {
	fu := &FileUpdate{}
	for _, fuw := range fuWiths {
		fuw(fu)
	}
	return *fu
}

func FuwPath(path string) FUWith {
	return func(fu *FileUpdate) {
		fu.Path = path
	}
}

func FuwStatus(status FileStatus) FUWith {
	return func(fu *FileUpdate) {
		fu.Status = status
	}
}

func FuwType(t string) FUWith {
	return func(fu *FileUpdate) {
		fu.Type = t
	}
}

func FuwProgress(progress float64) FUWith {
	return func(fu *FileUpdate) {
		fu.Progress = progress
	}
}

func FuwPeers(peers int) FUWith {
	return func(fu *FileUpdate) {
		fu.Peers = peers
	}
}

type FUWith func(fu *FileUpdate)

func FuwCid(cid string) FUWith {
	return func(fu *FileUpdate) {
		fu.CID = cid
	}
}

func FuwName(name string) FUWith {
	return func(fu *FileUpdate) {
		fu.Name = name
	}
}
