package ledger

import (
	"encoding/json"
	"fmt"
	"strings"

	aqb "github.com/umsatz/go-aqbanking"
)

func Parse(ts []aqb.Transaction) string {
	var text string
	for _, t := range ts {

		date := t.Date.Format("2006/01/02")
		credit := fmt.Sprintf("%f %s", t.Total, t.TotalCurrency)
		debit := fmt.Sprintf("%f %s", -t.Total, t.TotalCurrency)

		jsonString, _ := json.Marshal(t)

		text += fmt.Sprintf(
			";%s\n%s %s\n     %s  %s\n     %s  %s\n\n",
			jsonString,
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
