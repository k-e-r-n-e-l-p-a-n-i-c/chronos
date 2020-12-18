package cmd

import (
	"fmt"
	"os"

	"github.com/arunprasadmudaliar/chronos/pkg/controller"
	"github.com/spf13/cobra"
)

var kubeconfig string

var rootCmd = &cobra.Command{
	Use:   "chronos",
	Short: "Kubernetes event collector and notifier",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("Root command")
		config, _ := cmd.Flags().GetString("kubeconfig")
		controller.Start(config)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "path to kubeconfig file")
}
