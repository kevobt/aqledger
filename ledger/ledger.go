package ledger

import (
	"fmt"
	"strings"

	aqb "github.com/umsatz/go-aqbanking"
)

func ParseLedger(ts []aqb.Transaction) string {
	var text string
	for _, t := range ts {

		date := t.Date.Format("2006/01/02")
		credit := fmt.Sprintf("%f %s", t.Total, t.TotalCurrency)
		debit := fmt.Sprintf("%f %s", -t.Total, t.TotalCurrency)

		text += fmt.Sprintf(
			"%s %s\n     %s  %s\n     %s  %s\n\n",
			date,
			strings.Join(t.PurposeList, " "),
			"Assets",
			credit,
			"Expenses",
			debit,
		)
	}
	return text
}
