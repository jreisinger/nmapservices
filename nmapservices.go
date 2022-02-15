package nmapservices

import (
	"os"
)

var nmapServicesUrl = "https://raw.githubusercontent.com/nmap/nmap/master/nmap-services"
var nmapServicesFile = "/var/tmp/nmap-services"

func Top(n int) (Services, error) {
	if err := updateFile(nmapServicesFile, nmapServicesUrl); err != nil {
		return nil, err
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

	// Sanity checks.
	switch {
	case n > len(services)-1:
		n = len(services)
	case n < 0:
		n = 0
	}

	return services[:n], nil
}
