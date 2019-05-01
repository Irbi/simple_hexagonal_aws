package usecases

import (
	"fmt"
	"os"
	"test_aws/domain"
)

type File struct {
	Version int
	Name string
	Content []byte
}

type LocalInteractor struct {
	Repository domain.LocalFileRepository
}

type RemoteInteractor struct {
	Repository domain.RemoteFileRepository
}

type StorageInterfactor struct {
	Repository domain.StorageRepository
}

type FileInteractor struct {
	Remote *RemoteInteractor
	Local *LocalInteractor
	Storage *StorageInterfactor
}

func (itr *FileInteractor) SyncItems() (items []File, err error) {
	rFiles, err := itr.Remote.Repository.GetFiles()

	if err != nil {
		return nil, err
	}

	items = make([]File, 0)
	for _, rFile := range rFiles {
		lFile, err := itr.Local.Repository.FindByName(rFile.Name)
		if err != nil {
			return nil, err
		}

		isSynced := lFile.IsChecksumActual(rFile.Etag)
		if !isSynced {
			destPath := fmt.Sprintf("/tmp/test_aws/%s", lFile.Name)
			dest, err := os.Create(destPath)
			if err != nil {
				fmt.Printf("Unable to open file %q, %v", err)
				return nil, err
			}
			defer dest.Close()

			if (err != nil) {
				return nil, err
			}
			_, err = itr.Remote.Repository.Download(dest, lFile)
			if (err != nil) {
				return nil, err
			}
			err = itr.Storage.Repository.Store(lFile.Name, destPath)
			if (err != nil) {
				return nil, err
			}
			err = os.Remove(destPath)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			lFile.GetNextVersion()
			lFile.Checksum = rFile.Etag

			_, err = itr.Local.Repository.Store(lFile)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		}

		file := File{Name: lFile.Name, Version: lFile.Version, Content:lFile.Content}

		items = append(items, file)
	}

	return items, nil
}

func (itr *FileInteractor) GetFile(hash string, version int) (file File, err error) {

	file = File{}

	item, err := itr.Local.Repository.FindById(hash)
	if err != nil {
		return file, err
	}

	if item.Name != "" {
		content, err := itr.Storage.Repository.GetContent(item.Name)
		if err != nil {
			return file, err
		}
		file = File{Name: item.Name, Version: item.Version, Content: content}
	}

	return file, nil
}
