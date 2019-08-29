package ledger

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nikunjy/rules/parser"

	aqb "github.com/umsatz/go-aqbanking"
)

func Parse(ts []aqb.Transaction) string {
	var rules RulesCollection
	rules.ReadFromFile("rules")

	var text string
	from := "Assets"
	to := "Expenses"
	for _, t := range ts {
		for _, rule := range rules.Rules {
			ev, err := parser.NewEvaluator(rule.String)
			if err != nil {
				fmt.Printf("%v", err)
			}
			ans, err := ev.Process(transactionToMap(t))
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
	return text
}

func transactionToMap(t aqb.Transaction) map[string]interface{} {
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
