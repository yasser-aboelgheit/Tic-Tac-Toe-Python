/*
Copyright © 2024 George Messiha <georgemessiha22@gmail.com>
*/
package main

import (
	"log"

	"gitlab.com/startupbuilder/uda/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		log.Panicf("error: %v", err)
	}
}
