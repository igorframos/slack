package message

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

const (
	SenderTypeBot  SenderType = "bot"
	SenderTypeUser SenderType = "user"
)

type SenderType string

type Message struct {
	Channel      string `json:"channel"`
	Text         string `json:"text"`
	ParentThread string `json:"thread_ts,omitempty"`
	BotId        string `json:"bot_id,omitempty"`
	Type         string `json:"type,omitempty"`
	User         string `json:"user,omitempty"`
	Thread       string `json:"ts,omitempty"`
	Team         string `json:"team,omitempty"`
}

func createSendMessageRequest(senderType SenderType, m *Message) (*fasthttp.Request, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to create Send Message request")
	}

	req := fasthttp.AcquireRequest()
	req.SetBody(data)
	req.SetRequestURI("https://slack.com/api/chat.postMessage")
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json; charset=utf-8")

	if senderType != SenderTypeUser {
		return nil, errors.Errorf("Sender type %v is not supported yet", senderType)
	}

	// TODO: Use viper config to read token.
	token := os.Getenv("SLACK_USER_TOKEN")
	if senderType == SenderTypeBot {
		token = os.Getenv("SLACK_BOT_TOKEN")
	}
	req.Header.Add("Authorization", "Bearer "+token)

	return req, nil
}

func SendMessageAsUser(m *Message) (*Message, error) {
	req, err := createSendMessageRequest(SenderTypeUser, m)
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := fasthttp.Do(req, res); err != nil {
		return nil, errors.Wrapf(err, "Request to Slack API failed")
	}

	return readPostMessageResponse(res)
}
