package nmapservices

import "testing"

func TestTop(t *testing.T) {
	testcases := []struct {
		n, want int
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{10, 10},
		{10, 10},
		// grep -v '^#' /usr/local/share/nmap/nmap-services | wc -l
		{27425, 27425}, // max
		{27426, 27425},
		{100000, 27425},
	}

	services, err := Get()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testcases {
		got := len(services.Top(tc.n))
		if got != tc.want {
			t.Fatalf("got %d, wanted %d", got, tc.want)
		}
	}
}

func TestTcp(t *testing.T) {
	services, err := Get()
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range services.Tcp() {
		if s.Protocol != "tcp" {
			t.Fatalf("got %s, wanted tcp", s.Protocol)
		}
	}
}

func TestUdp(t *testing.T) {
	services, err := Get()
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range services.Udp() {
		if s.Protocol != "udp" {
			t.Fatalf("got %s, wanted udp", s.Protocol)
		}
	}
}

func TestTcpTop(t *testing.T) {
	testcases := []struct {
		n, want int
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{10, 10},
		// grep -v '^#' /usr/local/share/nmap/nmap-services | egrep '\s+\d+/tcp\s+' | wc -l
		{8351, 8351}, // max
		{8352, 8351},
		{100000, 8351},
	}

	services, err := Get()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testcases {
		got := len(services.Tcp().Top(tc.n))
		if got != tc.want {
			t.Fatalf("got %d, wanted %d", got, tc.want)
		}
	}
}

func TestUdpTop(t *testing.T) {
	testcases := []struct {
		n, want int
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{10, 10},
		// grep -v '^#' /usr/local/share/nmap/nmap-services | egrep '\s+\d+/udp\s+' | wc -l
		{19022, 19022}, // max
		{19023, 19022}, // max
		{100000, 19022},
	}

	services, err := Get()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testcases {
		got := len(services.Udp().Top(tc.n))
		if got != tc.want {
			t.Fatalf("got %d, wanted %d", got, tc.want)
		}
	}
}

func TestTopTcp(t *testing.T) {
	testcases := []struct {
		n, want int
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{10, 1},
		{100000, 8351},
	}

	services, err := Get()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testcases {
		got := len(services.Top(tc.n).Tcp())
		if got != tc.want {
			t.Fatalf("got %d, wanted %d", got, tc.want)
		}
	}
}

func TestTopUdp(t *testing.T) {
	testcases := []struct {
		n, want int
	}{
		{-1, 0},
		{0, 0},
		{1, 0},
		{10, 9},
		{100000, 19022},
	}

	services, err := Get()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testcases {
		got := len(services.Top(tc.n).Udp())
		if got != tc.want {
			t.Fatalf("got %d, wanted %d", got, tc.want)
		}
	}
}
