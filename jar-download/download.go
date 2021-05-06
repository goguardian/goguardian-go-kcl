package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
)

type downloader struct {
	maxRetries int
	backoff    time.Duration
	httpClient *http.Client

	mavenBaseURL string
	packages     []mavenPackageInfo
}

func getDownloader() *downloader {
	return &downloader{
		maxRetries: 3,
		backoff:    500 * time.Millisecond,
		httpClient: &http.Client{},

		packages:     mavenPackages,
		mavenBaseURL: mavenBaseHTTPURL,
	}
}

// Download fetches each jar package from Maven and saves it in the specified
// dstPath.
func (d *downloader) download(dstPath string) error {
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		err := os.Mkdir(dstPath, 0755)
		if err != nil {
			return errors.Wrap(err, "failed to make jar directory")
		}
	}

	for _, pkg := range d.packages {
		filename := path.Join(dstPath, pkg.Name())

		if err := d.downloadFileWithRetry(pkg.URL(d.mavenBaseURL), filename); err != nil {
			return err
		}
	}
	return nil
}

func (d *downloader) downloadFileWithRetry(src string, dst string) error {
	// don't download if dst already xists
	if _, err := os.Stat(dst); err == nil {
		return nil
	}
	var err error
	backoff := d.backoff

	for i := 0; i < d.maxRetries; i++ {
		fmt.Printf("Downloading %s to %s\n", src, dst)
		err = d.downloadFile(src, dst)
		if err == nil {
			break
		}

		time.Sleep(backoff)
		backoff *= 2 // exponentially backoff
	}

	if err != nil {
		return err
	}

	return nil
}

func (d *downloader) downloadFile(src string, dst string) error {
	resp, err := d.httpClient.Get(src)
	if err != nil {
		return errors.Wrap(err, "failed to download file")
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("non-2XX status code '%d'\n%s\n", resp.StatusCode, string(body))
	}

	out, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "failed to create destination file: %s", dst)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to write data to file")
	}

	return nil
}
