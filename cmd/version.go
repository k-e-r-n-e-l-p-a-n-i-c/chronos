package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "chronos version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Chronos version 1.0.0")
	},
}


