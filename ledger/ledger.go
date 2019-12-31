package ledger

import (
	"bufio"
	"encoding/json"
	"io"
	"regexp"
)

// A Ledger describes the contents of a ledger file filled with transaction
// information
type Ledger struct {
	Transactions Transactions
}

// AppendTransactions writes (appends) new transaction entries. The parameter rw
// defines where to write the transactions
func AppendTransactions(rw io.ReadWriter, ts []Transaction, rs []Rule, strict bool) error {
	ledger, err := Read(rw)
	if err != nil {
		return err
	}

	b, err := ParseTransactions(ledger.Transactions.Distinct(ts), rs, strict)
	if err != nil {
		return err
	}

	rw.Write(b)

	return nil
}

// Read reads aqledger transactions from an ledger stream
func Read(r io.Reader) (file Ledger, err error) {
	s := bufio.NewScanner(r)
	s.Split(scanTransactionLine)
	for s.Scan() {
		var t Transaction
		err = json.Unmarshal(s.Bytes(), &t)
		if err != nil {
			return
		}
		file.Transactions = append(file.Transactions, t)
	}

	return
}

// scanTransactionLine is a split method that looks for lines that contain
// transaction information
func scanTransactionLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Complile regular expression that expresses a transaction entry in a byte
	// source
	reg, err := regexp.Compile(";{.*}")
	if err != nil {
		return
	}

	// Find a transaction in data
	loc := reg.FindIndex(data)
	if loc != nil {
		return loc[1], data[loc[0]+1 : loc[1]], nil
	}

	// No entry found. Return
	return 0, nil, nil
}
