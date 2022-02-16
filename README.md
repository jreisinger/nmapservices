Sample use

```go
package main

import (
	"fmt"
	"log"

	"github.com/jreisinger/nmapservices"
)

func main() {
	services, err := nmapservices.Get()
	if err != nil {
		log.Fatal(err)
	}

	// Top 10 TCP ports you can find open on the Internet.
	for i, s := range services.Tcp().Top(10) {
		fmt.Printf("%2d. %5d # %s\n", i+1, s.Port, s.Name)
	}
}
```

Output

```
 1.    80 # http
 2.    23 # telnet
 3.   443 # https
 4.    21 # ftp
 5.    22 # ssh
 6.    25 # smtp
 7.  3389 # ms-wbt-server
 8.   110 # pop3
 9.   445 # microsoft-ds
10.   139 # netbios-ssn
```

See https://pkg.go.dev/github.com/jreisinger/nmapservices for more.
