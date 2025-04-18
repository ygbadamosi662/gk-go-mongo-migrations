package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "MongoDB migration manager",
}

func main() {
	// Register generate command
	var generateCmd = &cobra.Command{
		Use:   "generate [description]",
		Short: "Generate a new migration file",
		Args:  cobra.ExactArgs(1), // Expect exactly one argument (the description)
		Run: func(cmd *cobra.Command, args []string) {
			Generate(args[0])
		},
	}

	// Register init command
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the migration system (create config and migrations folder)",
		Run: func(cmd *cobra.Command, args []string) {
			InitDB() // Directly call InitDB from the main package
		},
	}

	// Add commands to root
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(initCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
