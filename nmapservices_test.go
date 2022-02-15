package nmapservices

import "testing"

func TestTop(t *testing.T) {
	testcases := []struct {
		want, get int
	}{
		{-1, 0},
		{0, 0},
		{10, 10},
		{27425, 27425}, // max
		{50000, 27425},
	}
	for _, tc := range testcases {
		services, err := Top(tc.want)
		if err != nil {
			t.Fatal(err)
		}
		if len(services) != tc.get {
			t.Fatalf("got %d, wanted %d", len(services), tc.get)
		}
	}
}

func TestTopTcp(t *testing.T) {
	services, err := TopTcp(27425)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range services {
		if s.Protocol != "tcp" {
			t.Fatalf("wanted tcp, got %s", s.Protocol)
		}
	}
}

func TestTopUdp(t *testing.T) {
	services, err := TopUdp(27425)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range services {
		if s.Protocol != "udp" {
			t.Fatalf("wanted tcp, got %s", s.Protocol)
		}
	}
}
