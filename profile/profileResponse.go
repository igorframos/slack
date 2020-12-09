package profile

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

// profileResponse is a helper to parse the data Slack returns to us.
type profileResponse struct {
	OK        bool     `json:"ok"`
	Profile   *Profile `json:"profile"`
	ErrorCode string   `json:"error,omitempty"`
}

func readProfileResponse(res *fasthttp.Response) (*Profile, error) {
	var result profileResponse
	if err := json.Unmarshal(res.Body(), &result); err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal profile information")
	}

	if !result.OK {
		return nil, errors.Errorf("Slack API reported an error: %v", result.ErrorCode)
	}

	if result.Profile == nil {
		return nil, errors.Errorf("Slack API returned an empty profile")
	}

	return result.Profile, nil
}
