package interfaces

type StorageHandler interface {
	Content (name string) (content []byte, err error)
	Store (name string, from string) error
}

type StorageRepo struct {
	handler  StorageHandler
}

func NewStorageRepo(handler StorageHandler) *StorageRepo {
	repo := new(StorageRepo)
	repo.handler = handler
	return repo
}

func (repo *StorageRepo) GetContent (name string) (content []byte, err error) {
	content, err = repo.handler.Content(name)
	if err != nil {
		return nil, err
	}

 	return content, nil
}

func (repo *StorageRepo) Store (name string, tmpPath string) (err error) {
	err = repo.handler.Store(name, tmpPath)

	return err
}