package cmd

import (
	"dev_scripts/cmd/forward"
	"dev_scripts/cmd/server"
	"dev_scripts/usecase"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewCmd() *cobra.Command {
	var license = "BSD 4-clause"
	var rootCmd = &cobra.Command{
		Use:   "cmd",
		Short: "All-in-one commands",
		Long:  "A collection of commands to do anything with less effort, basically",
	}

	env, err := usecase.LoadEnvFromFile()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	cmdServer := cmdserver.NewCmdServer(env)
	cmdForward := cmdforward.NewCmdForward(env)
	rootCmd.PersistentFlags().StringP("author", "a", "TomCrusade", "")
	rootCmd.PersistentFlags().StringVarP(&license, "license", "l", "", "")
	rootCmd.AddCommand(cmdServer.CobraCmd)
	rootCmd.AddCommand(cmdForward.CobraCmd)

	return rootCmd
}
