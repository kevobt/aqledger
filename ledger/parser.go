package ledger

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/nikunjy/rules/parser"
)

func ParseTransactions(ts []Transaction, rules []Rule, strict bool) ([]byte, error) {
	var text string
	uncategorized := []Transaction{}

	for _, t := range ts {
		// Remove \n from purpose
		re := regexp.MustCompile("\n")
		t.Purpose = re.ReplaceAllString(t.Purpose, "")

		from := "Assets"
		to := "Expenses"

		categorized := false
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
				categorized = true
				break
			}
		}

		if strict && !categorized {
			uncategorized = append(uncategorized, t)
			continue
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

	if len(uncategorized) > 0 {
		fmt.Println("Strict mode enabled. The following transactions could not be categorized:")

		for _, t := range uncategorized {
			data, err := json.MarshalIndent(t, "", "  ")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", data)
		}
	}

	return []byte(text), nil
}
