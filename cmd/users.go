package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	aqb "github.com/umsatz/go-aqbanking"
	"github.com/umsatz/go-aqbanking/examples"
)

var usersCmd = &cobra.Command{
	Use:   "users",
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

		for _, usr := range userCollection.Users {
			fmt.Printf("%s\n", usr.Name)
		}
	},
}
