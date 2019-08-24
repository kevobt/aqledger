package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kevobt/aqledger/ledger"
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
			fmt.Printf("%s", ledger.Parse(transactions))
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
