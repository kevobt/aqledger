package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kevobt/aqledger/ledger"
	"github.com/spf13/cobra"
)

var output string
var rules string

var rootCmd = &cobra.Command{
	Use:   "aqledger",
	Short: "aqledger is a tool to get transactions using HBCI and convertig them into ledger",
	Long:  "aqledger is a tool to get transactions using HBCI and convertig them into ledger",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You have to provide transactions")
			os.Exit(1)
		}

		var transactions ledger.Transactions

		err := json.Unmarshal([]byte(args[0]), &transactions)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if output != "" {
			rs := []ledger.Rule{}
			if rules != "" {
				// Read rules from file
				file, err := os.Open(rules)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer file.Close()
				rs, err = ledger.ReadRules(file)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			writeTransactionsToFile(output, transactions, rs)
		} else {
			printTransactions(transactions)
		}
	},
}

func init() {
	//rootCmd.AddCommand(usersCmd)
	//rootCmd.AddCommand(accountsCmd)
	//rootCmd.Args
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "output file")
	rootCmd.Flags().StringVarP(&rules, "rules", "r", "", "rules file")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func writeTransactionsToFile(filename string, ts ledger.Transactions, rs []ledger.Rule) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		ioutil.WriteFile(filename, []byte(""), 0644)
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	err = ledger.AppendTransactions(f, ts, rs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printTransactions(ts ledger.Transactions) {
	data, err := json.Marshal(ts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s", data)
}
