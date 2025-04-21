package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gk",
	Short: "MongoDB migration manager",
}

func main() {
	var generateCmd = &cobra.Command{
		Use:   "generate [description]",
		Short: "Generate a new migration file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Generate(args[0])
		},
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the migration system (create config and migrations folder)",
		Run: func(cmd *cobra.Command, args []string) {
			InitDB()
		},
	}

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
