package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {

}

var rootCmd = &cobra.Command{
	Use:   "aqledger",
	Short: "aqledger is a tool to get transactions using HBCI and convertig them into ledger",
	Long:  "aqledger is a tool to get transactions using HBCI and convertig them into ledger",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Hello World")
	},
}

// Execute runs the root commande
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
