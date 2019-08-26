package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/kevobt/aqledger/ledger"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	aqb "github.com/umsatz/go-aqbanking"
	"github.com/umsatz/go-aqbanking/examples"
)

var output string
var account string

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches bank transactions using HBCI",
	Long:  "Fetches bank transactions using HBCI",
	Run: func(cmd *cobra.Command, args []string) {
		aq, err := aqb.DefaultAQBanking()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer aq.Free()

		for _, pin := range examples.LoadPins("pins.json") {
			aq.RegisterPin(pin)
		}

		accountCollection, err := aq.Accounts()
		if err != nil {
			log.Fatalf("unable to list accounts: %v", err)
		}
		accountCollection = filterAccounts(accountCollection, func(a aqb.Account) bool {
			return a.Name == account
		})

		for _, account := range accountCollection {
			transactions, _ := aq.Transactions(&account, nil, nil)
			if output != "" {
				writeFile(output, transactions)
			} else {
				data, err := json.Marshal(transactions)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Printf("%s", data)
			}
		}
	},
}

func init() {
	fetchCmd.Flags().StringVarP(
		&account,
		"account",
		"a",
		"",
		"Transactions will be fetched from this account",
	)
}

func filterAccounts(as aqb.AccountCollection, f func(a aqb.Account) bool) aqb.AccountCollection {
	var asm aqb.AccountCollection
	for _, a := range as {
		if f(a) {
			asm = append(asm, a)
		}
	}
	return asm
}

func writeFile(filename string, newTransactions []aqb.Transaction) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		ioutil.WriteFile(filename, []byte(""), 0644)
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reg, err := regexp.Compile(";{.*}")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	res := reg.FindAll(b, -1)

	var old []aqb.Transaction
	var ts []aqb.Transaction

	for _, b := range res {
		var t aqb.Transaction
		// Skip ';' comment character
		err := json.Unmarshal(b[1:], &t)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		old = append(old, t)
	}

	for _, t := range newTransactions {
		equal := false
		for _, oldT := range old {
			if cmp.Equal(t, oldT) {
				equal = true
				break
			}
		}
		if !equal {
			ts = append(ts, t)
		}
	}

	ledgerString := ledger.Parse(ts)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	f.WriteString(ledgerString)
}
