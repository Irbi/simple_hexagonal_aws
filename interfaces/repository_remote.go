package interfaces

import (
	"os"
	"test_aws/domain"
	"time"
)

type RemoteFile interface {
	Key() string
	LastMod() time.Time
	Size() int64
	Etag() string
}

type RemoteDbHandler interface {
	Scan() (files []RemoteFile, err error)
	Download(dest *os.File, item string) (size int64, err error)
}

type RemoteFileRepo struct {
	dbHandler  RemoteDbHandler
}

func NewRemoteFileRepo(dbHandler RemoteDbHandler) *RemoteFileRepo {
	repo := new(RemoteFileRepo)
	repo.dbHandler = dbHandler
	return repo
}

func (repo *RemoteFileRepo) GetFiles() (files []domain.RemoteFile, err error) {
	items, err := repo.dbHandler.Scan()

	if err != nil {
		return nil, err
	}

	files = make([]domain.RemoteFile, len(items))
	for i, file := range items {
		eTag := file.Etag()
		if eTag[0] == '"' && eTag[len(eTag)-1] == '"' {
			eTag = eTag[1 : len(eTag)-1]
		}
		files[i] = domain.RemoteFile{Name:file.Key(), LastMod:file.LastMod(), Etag:eTag, Size:file.Size()}
	}

	return files, nil
}

func (repo *RemoteFileRepo) Download(fDest *os.File, file domain.File) (size int64, err error) {
	size, err = repo.dbHandler.Download(fDest, file.Name)
	return size, err
}

