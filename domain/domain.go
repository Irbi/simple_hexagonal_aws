package domain

import (
	guuid "github.com/google/uuid"
	"os"
	"time"
)

type LocalFileRepository interface {
	FindById(id string) (file File, err error)
	FindByName(name string) (file File, err error)
	Store(File) (result bool, err error)
}

type RemoteFileRepository interface {
	GetFiles() (files []RemoteFile, err error)
	Download(fDest *os.File, file File) (size int64, err error)
}

type StorageRepository interface {
	GetContent (name string) (content []byte, err error)
	Store (name string, tmpPath string) error
}

type StatRepository interface {
	AddValid() error
	AddOut() error
	Stat() (stat Stat, err error)
}

type File struct {
	ID string
	Name string
	Checksum string
	Version int
	Content []byte
}

type RemoteFile struct {
	Name string
	Size int64
	LastMod time.Time
	Etag string
}

type Stat struct {
	Total int
	Valid int
	Outdated int
}

func (file *File) IsChecksumActual(compareChecksum string) bool {
	return file.Checksum == compareChecksum
}

func (file *File) GenID() {
	if file.ID == "" {
		file.ID = guuid.New().String()
	}
}

func (file *File) GetNextVersion() {
	if file.ID != "" {
		file.Version++
	} else {
		file.Version = 1
	}
}

