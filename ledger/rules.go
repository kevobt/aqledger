package ledger

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

type Rule struct {
	String string
	From   string
	To     string
}

type RulesCollection struct {
	Rules []Rule
}

func (r *RulesCollection) ReadFromFile(filename string) error {
	f, _ := os.Open(filename)
	defer f.Close()

	scn := bufio.NewScanner(f)
	scn.Split(scanEmptyLine)

	var lines [][]byte

	for scn.Scan() {
		lines = append(lines, scn.Bytes())
	}

	for _, line := range lines {
		s := bufio.NewScanner(bytes.NewReader(line))
		var l []string
		for s.Scan() {
			l = append(l, s.Text())
		}

		r.Rules = append(r.Rules, Rule{
			String: l[0],
			From:   l[1],
			To:     l[2],
		})
	}

	return nil
}

func scanEmptyLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := strings.Index(string(data), "\n\n"); i >= 0 {
		return i + 2, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}
	return
}
