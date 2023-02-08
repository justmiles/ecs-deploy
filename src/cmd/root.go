package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	lambdaName   string
	debugEnabled bool
)

var rootCmd = &cobra.Command{
	Use:     "ecs-deploy",
	Short:   "Deploy something",
	Long:    `A fast and flexible tool to deploy to Amazon Web Service's Elastic Container Service`,
	Version: "0.7.1",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugEnabled, "debug", "d", false, "Enable debug logging")
}
