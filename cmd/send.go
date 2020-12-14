package cmd

import (
	"fmt"
	"strings"

	"github.com/igorframos/slack/message"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends a message as the user",
	Long:  `Sends a message as the user.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			return
		}

		channel, err := cmd.Flags().GetString("channel")
		if err != nil {
			fmt.Printf("Failed to get channel flag: %v\n", err)
			return
		}

		if !strings.HasPrefix(channel, "@") && !strings.HasPrefix(channel, "#") {
			fmt.Printf("The channel name must start with '@' or '#'")
			return
		}

		msg, err := message.SendMessageAsUser(&message.Message{
			Channel: channel,
			Text:    args[0],
		})

		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}

		fmt.Printf("Sent message. Thread ID: %s\n", msg.Thread)
	},
}

func init() {
	msgCmd.AddCommand(sendCmd)
}
