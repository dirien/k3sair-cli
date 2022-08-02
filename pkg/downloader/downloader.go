package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

type AirGapFileDownloaderImpl interface {
	Download(base, binary string) (*AirGapFile, error)
}

type AirGapFile struct {
	Path string
}

type AirGapeFileDownloader struct{}

func (a AirGapeFileDownloader) Download(base, binary string) (*AirGapFile, error) {
	fmt.Printf("Download Air-Gap file %s\n", color.GreenString(binary))
	tmp := os.TempDir()

	p := filepath.FromSlash(fmt.Sprintf("%s/%s", tmp, binary))
	out, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	var transport http.RoundTripper = &http.Transport{
		DisableKeepAlives: true,
	}
	c := &http.Client{Transport: transport}

	resp, err := c.Get(fmt.Sprintf("%s/%s", base, binary))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Air-Gap file succesfully downloaded at %s\n", color.RedString(p))
	return &AirGapFile{Path: p}, nil
}
