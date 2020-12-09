package cmd

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/igorframos/slack/profile"
)

const (
	durationFlagName   = "duration"
	emojiFlagName      = "emoji"
	expirationFlagName = "expiration"
)

var (
	acceptedTimeFormats = []string{
		"2/1 3:04PM",
		"2/1 15:04",
		"3:04PM 2/1",
		"15:04 2/1",
		"3:04PM 2/1/2006",
		"15:04 2/1/2006",
		"2006/1/2 15:04",
		"2006/1/2 3:04PM",
		"2006-01-02 15:04",
		"2006-01-02 3:04PM",
		"2/1/2006 15:04",
		"2/1/2006 3:04PM",
		"3:04PM",
		"15:04",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
	}
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets the status of the user",
	Long:  `Sets the status of the user.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			return
		}

		emoji, err := cmd.Flags().GetString(emojiFlagName)
		if err != nil {
			fmt.Printf("Failed to get emoji flag: %v\n", err)
			return
		}

		expiration, err := getExpiration(cmd)
		if err != nil {
			fmt.Printf("Error determining expiration: %v\n", err)
			return
		}

		status := &profile.Status{
			Text:  args[0],
			Emoji: emoji,
		}
		if expiration > 0 {
			status.Expiration = expiration
		}

		status, err = profile.SetStatus(status)
		if err != nil {
			fmt.Printf("Failed to set status: %v\n", err)
			return
		}

		fmt.Printf("Status successfuly set.\n\n%v\n", status)
	},
}

func getExpiration(cmd *cobra.Command) (int64, error) {
	duration, err := cmd.Flags().GetDuration(durationFlagName)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to get `%s` flag", durationFlagName)
	}

	expirationStr, err := cmd.Flags().GetString(expirationFlagName)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to get `%s` flag", expirationFlagName)
	}

	expireAt, err := parseExpiration(expirationStr)
	if err != nil {
		return 0, errors.Wrapf(err, "Error in flag `%s`: cannot parse `%s` as a time", expirationFlagName, expirationStr)
	}

	var expiration time.Time
	if duration != 0 {
		expiration = time.Now().Add(duration)
	} else {
		expiration = expireAt
	}

	return expiration.Unix(), nil
}

// parseExpiration tries several different formats to parse a time. It tries to give some flexibility to the flag value.
func parseExpiration(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, nil
	}

	for _, layout := range acceptedTimeFormats {
		t, err := time.ParseInLocation(layout, raw, time.Local)
		if err == nil {
			return pickNextDate(t), nil
		}
	}

	return time.Time{}, errors.Errorf("Could not parse `%s` as a time", raw)
}

func pickNextDate(t time.Time) time.Time {
	y, m, d := t.Date()

	if y != 0 {
		return t
	}

	if y == 0 {
		y = time.Now().Year()
	}

	if t.Month() == 1 {
		m = time.Now().Month()
	}

	if t.Day() == 1 {
		d = time.Now().Day()
	}

	t = time.Date(y, m, d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())

	if t.Before(time.Now()) {
		t = t.AddDate(0, 0, 1)
	}

	return t
}

func init() {
	statusCmd.AddCommand(setCmd)

	setCmd.Flags().StringP(emojiFlagName, "e", "", "The emoji to use for the status, e.g. :cookie:.")
	setCmd.Flags().DurationP(durationFlagName, "d", 0, "An amount of time after which the status should be cleared. Duration must be described in hours (h), minutes (m), seconds (s) or smaller units.")
	setCmd.Flags().StringP(expirationFlagName, "x", "", "A time at which the status should be automatically cleared.")
}
