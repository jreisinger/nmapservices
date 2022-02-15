package nmapservices

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Service struct {
	Name      string
	PortProto string // e.g. 22/tcp
	Frequency float64
	Comment   string // optional
}

type Services []Service

func parseServiceFile(file *os.File) (Services, error) {
	var ss Services

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
		svc := Service{
			Name:      parts[0],
			PortProto: parts[1],
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

func (ss Services) print(top uint, noHeader bool) {
	sort.Sort(sort.Reverse(byFrequency(ss)))

	const format = "%v\t%v\t%v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	if !noHeader {
		fmt.Fprintf(tw, format, "Name", "Port/Proto", "Frequency", "Comment")
		fmt.Fprintf(tw, format, "-----", "---------", "---------", "-------")
	}
	for _, s := range ss[:top] {
		fmt.Fprintf(tw, format, s.Name, s.PortProto, s.Frequency, s.Comment)
	}
	tw.Flush() // calculate column widths and print table
}

type byFrequency Services

func (x byFrequency) Len() int           { return len(x) }
func (x byFrequency) Less(i, j int) bool { return x[i].Frequency < x[j].Frequency }
func (x byFrequency) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
