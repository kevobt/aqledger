package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kevobt/aqledger/ledger"
	"github.com/spf13/cobra"
	aqb "github.com/umsatz/go-aqbanking"
	"github.com/umsatz/go-aqbanking/examples"
)

var output string
var rules string

var rootCmd = &cobra.Command{
	Use:   "aqledger",
	Short: "aqledger is a tool to get transactions using HBCI and convertig them into ledger",
	Long:  "aqledger is a tool to get transactions using HBCI and convertig them into ledger",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You have to provide an account")
			fmt.Println("You can use the command  \"aqledger accounts\" to list your accounts ")
			os.Exit(1)
		}
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
			return a.Name == args[0]
		})

		for _, account := range accountCollection {
			transactions, _ := aq.Transactions(&account, nil, nil)
			var ts ledger.Transactions
			for _, t := range transactions {
				ts = append(ts, ledger.Transaction(t))
			}
			t := ledger.Transactions(ts)

			if output != "" {
				writeTransactionsToFile(output, t)
			} else {
				printTransactions(t)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(accountsCmd)
	rootCmd.AddCommand(transactionsCmd)
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "output file")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func filterAccounts(as aqb.AccountCollection, f func(a aqb.Account) bool) (asm aqb.AccountCollection) {
	for _, a := range as {
		if f(a) {
			asm = append(asm, a)
		}
	}
	return asm
}

func writeTransactionsToFile(filename string, ts ledger.Transactions) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		ioutil.WriteFile(filename, []byte(""), 0644)
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	ledger.AppendTransactions(f, ts)
}

func printTransactions(ts ledger.Transactions) {
	data, err := json.Marshal(ts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s", data)
}
