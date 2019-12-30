package ledger

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/nikunjy/rules/parser"
)

func ParseTransactions(ts []Transaction, rules []Rule) ([]byte, error) {
	var text string
	for _, t := range ts {
		// Remove \n from purpose
		re := regexp.MustCompile("\n")
		t.Purpose = re.ReplaceAllString(t.Purpose, "")

		from := "Assets"
		to := "Expenses"
		for _, rule := range rules {
			ev, err := parser.NewEvaluator(rule.String)
			if err != nil {
				fmt.Println("There might be some mistake in the rules definition")
				fmt.Println(err)
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

		credit := fmt.Sprintf("%0.2f %s", t.Total, t.TotalCurrency)

		jsonString, _ := json.Marshal(t)

		text += fmt.Sprintf(
			";%s\n%s %s\n     %s  %s\n     %s  %s\n\n",
			jsonString,
			t.Date,
			t.Purpose,
			from,
			credit,
			to,
			"",
		)
	}

	return []byte(text), nil
}
