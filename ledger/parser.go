package ledger

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nikunjy/rules/parser"
)

func ParseTransactions(ts []Transaction, rules []Rule) ([]byte, error) {
	var text string
	from := "Assets"
	to := "Expenses"
	for _, t := range ts {
		for _, rule := range rules {
			ev, err := parser.NewEvaluator(rule.String)
			if err != nil {
				fmt.Printf("%v", err)
			}
			ans, err := ev.Process(t.Map())
			if err != nil {
				fmt.Println(err)
			}
			if ans {
				from = rule.From
				to = rule.To
				break
			}
		}

		date := t.Date.Format("2006/01/02")
		credit := fmt.Sprintf("%f %s", -t.Total, t.TotalCurrency)
		debit := fmt.Sprintf("%f %s", t.Total, t.TotalCurrency)

		jsonString, _ := json.Marshal(t)

		text += fmt.Sprintf(
			";%s\n%s %s\n     %s  %s\n     %s  %s\n\n",
			jsonString,
			date,
			strings.Join(t.PurposeList, " "),
			from,
			credit,
			to,
			debit,
		)
	}
	return []byte(text), nil
}
