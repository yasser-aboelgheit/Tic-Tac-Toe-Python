/*
Copyright © 2024 George Messiha <georgemessiha22@gmail.com>
*/
package main

import (
	"log"

	"github.com/spf13/cobra"
	"gitlab.com/startupbuilder/uda/cmd"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "",
		Short: "start any service.",
		Long:  "this command should be replaced later with the app/root.go NewRootCmd directly.",
	}

	rootCmd.AddCommand(cmd.NewRootCmd())

	return rootCmd
}

func main() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		log.Panicf("error: %v", err)
	}
}
