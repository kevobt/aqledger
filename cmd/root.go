package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
	rootCmd.AddCommand(usersCmd)
	fetchCmd.LocalFlags().StringVarP(&output, "output", "o", "", "output file")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
