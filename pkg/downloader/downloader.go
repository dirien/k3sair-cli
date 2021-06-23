package downloader

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type AirGapFileDownloaderImpl interface {
	Download(base, binary string) (*AirGapFile, error)
}

type AirGapFile struct {
	Path string
}

type AirGapeFileDownloader struct {
}

func (a AirGapeFileDownloader) Download(base, binary string) (*AirGapFile, error) {
	fmt.Println(fmt.Sprintf("Download Air-Gap file %s", color.GreenString(binary)))
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

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
	fmt.Println(fmt.Sprintf("Air-Gap file succesfully downloaded at %s", color.RedString(p)))
	return &AirGapFile{Path: p}, nil
}
