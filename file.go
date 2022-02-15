package nmapservices

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// updateFile updates file from url if the file is older than a week. If file
// does not exist it downloads and creates it.
func updateFile(file, url string) error {
	f, err := os.Stat(file)

	if os.IsNotExist(err) {
		r, err := downloadFile(url)
		if err != nil {
			return err
		}
		if err := storeFile(file, r); err != nil {
			return err
		}

		return nil // don't check ModTime if file does not exist
	}

	if isOlderThanOneWeek(f.ModTime()) {
		r, err := downloadFile(url)
		if err != nil {
			return err
		}
		if err := storeFile(file, r); err != nil {
			return err
		}
	}

	return nil
}

func storeFile(outFilename string, r io.ReadCloser) error {
	defer r.Close() // let's close resp.Body

	outFile, err := os.Create(outFilename)
	if err != nil {
		return nil
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, r); err != nil {
		return err
	}

	return nil
}

func isOlderThanOneWeek(t time.Time) bool {
	return time.Since(t) > 7*24*time.Hour
}

func downloadFile(url string) (r io.ReadCloser, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Check the server response.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't download %v: %v", url, resp.Status)
	}

	return resp.Body, nil
}
