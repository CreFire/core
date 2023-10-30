package main

import (
	"fmt"
	"github.com/core/cmd/excel"
	"github.com/core/cmd/svn"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	Execute()
	//excel.Token()
}

var rootCmd = &cobra.Command{
	Use:     "gen",
	Short:   "root cmd",
	Version: "1.0.0",
}

func Execute() {
	rootCmd.AddCommand(svn.SvnUpCmd)
	rootCmd.AddCommand(excel.ExcelCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
