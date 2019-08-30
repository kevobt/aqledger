package ledger

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// Rule represents a condition that, if satisfied, indicates where money flows
type Rule struct {
	String string
	From   string
	To     string
}

// ReadRules reads rules from a reader. A rule consists of three lines.
// 1. Rule expression
// 2. Source account
// 3. Target account
// This means, if a transaction satisfies a condition the money flows from the
// source account to the target account. Multiple rules are separated by a blank
// line.
func ReadRules(r io.Reader) (rules []Rule, err error) {
	rulesScanner := bufio.NewScanner(r)
	rulesScanner.Split(scanEmptyLine)

	for rulesScanner.Scan() {
		linesScanner := bufio.NewScanner(bytes.NewReader(rulesScanner.Bytes()))
		var lines []string
		for linesScanner.Scan() {
			lines = append(lines, linesScanner.Text())
		}
		rules = append(rules, Rule{
			String: lines[0],
			From:   lines[1],
			To:     lines[2],
		})
	}

	return rules, nil
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
