package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/igorframos/slack/profile"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears the status of the user",
	Long:  `Clears the status of the user.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := profile.ClearStatus(); err != nil {
			fmt.Printf("Failed to clear status: %v\n", err)
			return
		}

		fmt.Printf("Status successfuly cleared.\n")
	},
}

func init() {
	statusCmd.AddCommand(clearCmd)
}
