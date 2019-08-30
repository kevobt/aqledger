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

var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "aqledger is a tool to get transactions using HBCI and convertig them into ledger",
	Long:  "aqledger is a tool to get transactions using HBCI and convertig them into ledger",
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
			var ts ledger.Transactions
			for _, t := range transactions {
				ts = append(ts, ledger.Transaction(t))
			}
			t := ledger.Transactions(ts)
			fmt.Printf("%s", t.String())
		}
	},
}

func init() {
	transactionsCmd.Flags().StringVarP(
		&account,
		"account",
		"a",
		"",
		"Transactions will be fetched from this account",
	)
}
