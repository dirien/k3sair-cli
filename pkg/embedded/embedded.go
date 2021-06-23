package embedded

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type EmbeddedFileLoaderImpl interface {
	LoadEmbeddedFile(content, tmpFilename string) (*EmbeddedFile, error)
}

type EmbeddedFile struct {
	Path string
}

type EmbeddedFileLoader struct {
}

func (e EmbeddedFileLoader) LoadEmbeddedFile(content, tmpFilename string) (*EmbeddedFile, error) {
	tmp, err := ioutil.TempDir("", "")
	p := filepath.FromSlash(fmt.Sprintf(tmpFilename, tmp))
	err = ioutil.WriteFile(p, []byte(content), 0644)
	if err != nil {
		return nil, err
	}
	return &EmbeddedFile{
		Path: p,
	}, err
}
