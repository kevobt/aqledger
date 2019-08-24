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

		userCollection, err := aq.Users()
		if err != nil {
			log.Fatalf("unable to list users: %v", err)
		}
		defer userCollection.Free()

		accountCollection, err := aq.Accounts()
		if err != nil {
			log.Fatalf("unable to list accounts: %v", err)
		}

		for _, account := range accountCollection {
			transactions, _ := aq.Transactions(&account, nil, nil)
			fmt.Printf("%s", ledger.ParseLedger(transactions))
		}
	},
}

var rootCmd = &cobra.Command{
	Use:   "aqledger",
	Short: "aqledger is a tool to get transactions using HBCI and convertig them into ledger",
	Long:  "aqledger is a tool to get transactions using HBCI and convertig them into ledger",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Hello World")
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.LocalFlags().StringVarP(&output, "output", "o", "", "output file")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
