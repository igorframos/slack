package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// msgCmd represents the msg command
var msgCmd = &cobra.Command{
	Use:   "msg",
	Short: "Allows some basic message manipulation",
	Long:  `Allows some basic message manipulation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("msg called")
	},
}

func init() {
	rootCmd.AddCommand(msgCmd)

	msgCmd.PersistentFlags().StringP("channel", "c", "", "The channel to use in the command. Must start with @ or #")

	msgCmd.PersistentFlags().StringP("thread", "t", "", "The thread to use in the command")
}
