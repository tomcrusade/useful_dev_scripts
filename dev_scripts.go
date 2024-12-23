package main

import (
	"dev_scripts/cmd"
	"fmt"
	"os"
)

func main() {
	rootCmd := cmd.NewCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
