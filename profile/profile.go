package profile

import (
	"os"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type Profile struct {
	*Status

	DisplayName string `json:"display_name,omitempty"`
	Email       string `json:"email,omitempty"`
	RealName    string `json:"real_name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	ImageURL    string `json:"image_original,omitempty"`
}

func createGetProfileRequest() *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://slack.com/api/users.profile.get")

	// TODO: Use viper config to read token.
	req.Header.Add("Authorization", "Bearer "+os.Getenv("SLACK_USER_TOKEN"))

	return req
}

func ReadProfile() (*Profile, error) {
	req := createGetProfileRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := fasthttp.Do(req, res); err != nil {
		return nil, errors.Wrapf(err, "Failed to retrieve profile information")
	}

	return readProfileResponse(res)
}
