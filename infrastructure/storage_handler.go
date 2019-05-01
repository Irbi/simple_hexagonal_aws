package infrastructure

import (
	"fmt"
	"io"
	"os"
)

type LocalFileHandler struct {
	Path string
}

func NewStorageHandler(path string) *LocalFileHandler {
	handler := new(LocalFileHandler)
	handler.Path = path

	return handler
}

func (handler *LocalFileHandler) Content (name string) (content []byte, err error) {
	file, err := os.Open(fmt.Sprintf("%s/%s", handler.Path, name))
	defer file.Close()
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func (handler *LocalFileHandler) Store (name string, from string) error {
	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(fmt.Sprintf("%s/%s", handler.Path, name))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)

	return err
}

