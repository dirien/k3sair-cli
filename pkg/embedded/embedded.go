package embedded

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileLoaderImpl interface {
	LoadEmbeddedFile(content, tmpFilename string) (*File, error)
}

type File struct {
	Path string
}

type FileLoader struct{}

func (e FileLoader) LoadEmbeddedFile(content, tmpFilename string) (*File, error) {
	tmp := os.TempDir()
	p := filepath.FromSlash(fmt.Sprintf(tmpFilename, tmp))
	err := os.WriteFile(p, []byte(content), 0o644)
	if err != nil {
		return nil, err
	}
	return &File{
		Path: p,
	}, err
}
