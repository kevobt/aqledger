package ledger

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

// Transaction represents a banking transaction indicating money flow from one
// account to another
type Transaction struct {
	Purpose             string  `json:"purpose"`
	Date                string  `json:"date"`
	ValutaDate          string  `json:"valutaDate"`
	Total               float32 `json:"total"`
	TotalCurrency       string  `json:"totalCurrency"`
	LocalBankCode       string  `json:"localBankCode"`
	LocalAccountNumber  string  `json:"localAccountNumber"`
	RemoteBankCode      string  `json:"remoteBankCode"`
	RemoteAccountNumber string  `json:"remoteAccountNumber"`
	RemoteName          string  `json:"remoteName"`
}

// Map transforms the transaction struct into a map. Each key represents a
// transactions property.
func (t Transaction) Map() map[string]interface{} {
	type s map[string]interface{}
	return s{
		"Purpose":             t.Purpose,
		"Date":                t.Date,
		"ValutaDate":          t.ValutaDate,
		"Total":               t.Total,
		"TotalCurrency":       t.TotalCurrency,
		"LocalBankCode":       t.LocalBankCode,
		"LocalAccountNumber":  t.LocalAccountNumber,
		"RemoteBankCode":      t.RemoteBankCode,
		"RemoteAccountNumber": t.RemoteAccountNumber,
		"RemoteName":          t.RemoteName,
	}
}

// String prints transaction in a human readable format
func (t Transaction) String() string {
	return fmt.Sprintf("%s %s\n%f %s", t.Date, t.Purpose, t.Total, t.TotalCurrency)
}

// Transactions describes a slice of transactions
type Transactions []Transaction

// Distinct compares a transaction list with another one and returns all
// transactions that are only in the others list.
func (ts Transactions) Distinct(others Transactions) (distinct []Transaction) {
	for _, o := range others {
		found := false
		for _, t := range ts {
			if cmp.Equal(o, t) {
				found = true
				break
			}
		}
		if !found {
			distinct = append(distinct, o)
		}
	}
	return
}

// String returns transactions in a human readable format
func (ts Transactions) String() (val string) {
	for _, t := range ts {
		val += fmt.Sprintf("%s\n\n", t.String())
	}
	return
}
