package cmd

import (
	"github.com/spf13/cobra"
)

var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Migration parent command",
	Run:   runMigrationCmd,
}

func init() {
	rootCmd.AddCommand(migrationCmd)
}

func runMigrationCmd(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		panic(err)
	}
}
