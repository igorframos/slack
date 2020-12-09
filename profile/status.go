package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

var (
	emptyStatus = &Status{
		Text:       "",
		Expiration: 0,
	}
)

type Status struct {
	Emoji      string `json:"status_emoji"`
	Text       string `json:"status_text"`
	Expiration int64  `json:"status_expiration"`
}

// setStatusRequest is a helper to marshal the request properly.
type setStatusRequest struct {
	Profile *Profile `json:"profile"`
}

func createSetStatusRequest(s *Status) (*fasthttp.Request, error) {
	data, err := json.Marshal(setStatusRequest{&Profile{Status: s}})
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to create Set Status request")
	}

	req := fasthttp.AcquireRequest()
	req.SetBody(data)
	req.SetRequestURI("https://slack.com/api/users.profile.set")
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json; charset=utf-8")

	// TODO: Use viper config to read token.
	req.Header.Add("Authorization", "Bearer "+os.Getenv("SLACK_USER_TOKEN"))

	return req, nil
}

func SetStatus(s *Status) (*Status, error) {
	req, err := createSetStatusRequest(s)
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := fasthttp.Do(req, res); err != nil {
		return nil, errors.Wrapf(err, "Request to Slack API failed")
	}

	profile, err := readProfileResponse(res)
	if err != nil {
		return nil, err
	}

	return profile.Status, nil
}

func ClearStatus() error {
	req, err := createSetStatusRequest(emptyStatus)
	if err != nil {
		return err
	}
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := fasthttp.Do(req, res); err != nil {
		return errors.Wrapf(err, "Failed to clear status")
	}

	_, err = readProfileResponse(res)

	return err
}

func (s *Status) String() string {
	expiration := "Status will not expire."
	if s.Expiration > 0 {
		expiration = fmt.Sprintf("Status set to expire at %v", time.Unix(s.Expiration, 0))
	}

	return fmt.Sprintf("(%s) %s\n\n%s", s.Emoji, s.Text, expiration)
}
