package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/igorframos/slack/profile"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Performs operations on your Slack status",
	Long: `Performs operations on your Slack status. The following operations are available:

set STATUS_TEXT [FLAGS]
  Allows setting your status. Example:

  slack status set 'Eating chocolate' --emoji ':chocolate_bar:'`,
	Run: func(cmd *cobra.Command, args []string) {
		prof, err := profile.ReadProfile()
		if err != nil {
			fmt.Printf("Failed to fetch Slack status: %v\n", err)
			return
		}

		if prof.Status == nil {
			fmt.Printf("Error fetching status information. Response from Slack does not contain expected status information.\n")
			return
		}

		if prof.Status.Text == "" && prof.Status.Emoji == "" {
			fmt.Printf("No status currently set.\n")
			return
		}

		fmt.Printf("%v\n", prof.Status)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
