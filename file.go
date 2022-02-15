package nmapservices

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// NmapServicesFileUrl is the online location of the nmap-services file.
var NmapServicesFileUrl = "https://raw.githubusercontent.com/nmap/nmap/master/nmap-services"

// NmapServicesFiles are typical filesystem locations of the nmap-services file.
var NmapServicesFiles = []string{
	"/usr/share/nmap/nmap-services",
	"/usr/local/share/nmap/nmap-services",
}

// GetServices extract Services from nmap-services file. First if tries
// NmapServicesFiles. If none is present present it downloads the file from
// NmapServicesFileUrl.
func GetServices() ([]Service, error) {
	var nmapServicesFile string

	for _, f := range NmapServicesFiles {
		if _, err := os.Open(f); err == nil {
			nmapServicesFile = f
			break
		}
	}

	if nmapServicesFile == "" {
		nmapServicesFile = "/var/tmp/nmap-services"
		if err := updateFile(nmapServicesFile, NmapServicesFileUrl); err != nil {
			return nil, err
		}
	}

	file, err := os.Open(nmapServicesFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	services, err := parseServiceFile(file)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func parseServiceFile(file *os.File) ([]Service, error) {
	var ss []Service

	input := bufio.NewScanner(file)
	ws := regexp.MustCompile(`\s+`)
	for input.Scan() {
		if strings.HasPrefix(input.Text(), "#") { // skip comments
			continue
		}
		parts := ws.Split(input.Text(), 4)
		freq, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		var comment string
		if len(parts) == 4 {
			comment = parts[3]
		}
		portProto := strings.Split(parts[1], "/")
		p, err := strconv.Atoi(portProto[0])
		if err != nil {
			return nil, err
		}
		svc := Service{
			Name:      parts[0],
			Port:      int16(p),
			Protocol:  portProto[1],
			Frequency: freq,
			Comment:   comment,
		}
		ss = append(ss, svc)
	}
	if err := input.Err(); err != nil {
		log.Fatal(err)
	}

	return ss, nil
}

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
