package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	aqb "github.com/umsatz/go-aqbanking"
	"github.com/umsatz/go-aqbanking/examples"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Lists aqbanking accounts",
	Long:  "Lists aqbanking accounts",
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

		for _, acc := range accountCollection {
			fmt.Printf("%s\n", acc.Name)
		}
	},
}
