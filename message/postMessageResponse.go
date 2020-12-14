package message

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type postMessageResponse struct {
	OK        bool     `json:"ok"`
	Channel   string   `json:"channel"`
	Thread    string   `json:"ts"`
	Message   *Message `json:"message,omitempty"`
	ErrorCode string   `json:"error,omitempty"`
}

func readPostMessageResponse(res *fasthttp.Response) (*Message, error) {
	var result postMessageResponse
	if err := json.Unmarshal(res.Body(), &result); err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal message information")
	}

	if !result.OK {
		return nil, errors.Errorf("Slack API reported an error: %v", result.ErrorCode)
	}

	if result.Message == nil {
		return nil, errors.Errorf("Slack API returned an empty message as response")
	}

	return result.Message, nil
}
