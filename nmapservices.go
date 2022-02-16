// Package nmapservices provides network services and the frequency they are
// found on the Internet. See https://nmap.org/book/nmap-services.html for more.
package nmapservices

import (
	"os"
	"sort"
)

// NmapServicesUrl is the online location of the nmap-services file.
var NmapServicesUrl = "https://raw.githubusercontent.com/nmap/nmap/master/nmap-services"

// NmapServicesFiles are typical filesystem locations of the nmap-services file.
var NmapServicesFiles = []string{
	"/usr/share/nmap/nmap-services",
	"/usr/local/share/nmap/nmap-services",
}

// Service represents a network service.
type Service struct {
	Name      string
	Port      int16  // e.g. 22
	Protocol  string // e.g. tcp
	Frequency float64
	Comment   string // optional
}

type Services []Service

// Get extracts Services from nmap-services file. First if tries
// NmapServicesFiles. If none is present locally it downloads the file from
// NmapServicesUrl.
func Get() (Services, error) {
	var nmapServicesFile string

	for _, f := range NmapServicesFiles {
		if _, err := os.Open(f); err == nil {
			nmapServicesFile = f
			break
		}
	}

	if nmapServicesFile == "" {
		nmapServicesFile = "/var/tmp/nmap-services"
		if err := updateFile(nmapServicesFile, NmapServicesUrl); err != nil {
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

// Top returns n services that are found most frequently on the Internet.
func (services Services) Top(n int) Services {
	return getTopN(services, n)
}

// Tcp returns services using TCP protocol.
func (services Services) Tcp() Services {
	var t Services
	for _, s := range services {
		if s.Protocol == "tcp" {
			t = append(t, s)
		}
	}
	services = t
	return services
}

// Udp returns services using UDP protocol.
func (services Services) Udp() Services {
	var t Services
	for _, s := range services {
		if s.Protocol == "udp" {
			t = append(t, s)
		}
	}
	services = t
	return services
}

func getTopN(services []Service, n int) []Service {
	// Sanity checks.
	switch {
	case n > len(services)-1:
		n = len(services)
	case n < 0:
		n = 0
	}

	sort.Sort(sort.Reverse(byFrequency(services)))

	return services[:n]
}

type byFrequency []Service

func (x byFrequency) Len() int           { return len(x) }
func (x byFrequency) Less(i, j int) bool { return x[i].Frequency < x[j].Frequency }
func (x byFrequency) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
