// Package nmapservices provides network services and the frequency they are
// found on the Internet. See https://nmap.org/book/nmap-services.html for more.
package nmapservices

import (
	"sort"
)

// Service represents a network service.
type Service struct {
	Name      string
	Port      int16  // e.g. 22
	Protocol  string // e.g. tcp
	Frequency float64
	Comment   string // optional
}

// Top returns n services that are found most frequently on the Internet.
func Top(n int) ([]Service, error) {
	services, err := GetServices()
	if err != nil {
		return nil, err
	}

	return getTopN(services, n), nil
}

// TopTcp returns n services using TCP protocol that are found most frequently
// on the Internet.
func TopTcp(n int) ([]Service, error) {
	services, err := GetServices()
	if err != nil {
		return nil, err
	}

	var t []Service
	for _, s := range services {
		if s.Protocol == "tcp" {
			t = append(t, s)
		}
	}
	services = t

	return getTopN(services, n), nil
}

// TopUdp returns n services using UDP protocol that are found most frequently
// on the Internet.
func TopUdp(n int) ([]Service, error) {
	services, err := GetServices()
	if err != nil {
		return nil, err
	}

	var t []Service
	for _, s := range services {
		if s.Protocol == "udp" {
			t = append(t, s)
		}
	}
	services = t

	return getTopN(services, n), nil
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
