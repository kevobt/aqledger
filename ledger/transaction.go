package ledger

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	aqb "github.com/umsatz/go-aqbanking"
)

// Transaction represents a banking transaction indicating money flow from one
// account to another
type Transaction aqb.Transaction

// Map transforms the transaction struct into a map. Each key represents a
// transactions property.
func (t Transaction) Map() map[string]interface{} {
	type s map[string]interface{}
	return s{
		"Purpose":             t.Purpose,
		"PurposeList":         t.PurposeList,
		"Text":                t.Text,
		"Status":              t.Status,
		"Date":                t.Date,
		"ValutaDate":          t.ValutaDate,
		"CustomerReference":   t.CustomerReference,
		"EndToEndReference":   t.EndToEndReference,
		"Total":               t.Total,
		"TotalCurrency":       t.TotalCurrency,
		"Fee":                 t.Fee,
		"FeeCurrency":         t.FeeCurrency,
		"MandateID":           t.MandateID,
		"BandReference":       t.BandReference,
		"LocalBankCode":       t.LocalBankCode,
		"LocalAccountNumber":  t.LocalAccountNumber,
		"LocalIBAN":           t.LocalIBAN,
		"LocalBIC":            t.LocalBIC,
		"LocalName":           t.LocalName,
		"RemoteBankCode":      t.RemoteBankCode,
		"RemoteAccountNumber": t.RemoteAccountNumber,
		"RemoteIBAN":          t.RemoteIBAN,
		"RemoteBIC":           t.RemoteBIC,
		"RemoteName":          t.RemoteName,
		"RemoteNameList":      t.RemoteNameList,
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
		for _, t := range ts {
			if cmp.Equal(o, t) {
				continue
			}
		}
		distinct = append(distinct, o)
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
