package services

import "maildefender/engine/internal/configuration"

type SendMailIn struct {
	ReplyTo     string   `json:"replyTo"`
	ThreadTopic string   `json:"threadTopic"`
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
}

func SendEmail(in SendMailIn) (int, error) {
	return doReq(
		configuration.NotifierBaseEndpoint(),
		request{
			reqType:  post,
			endpoint: "/v1/notifier/email",
			body:     in,
		}, nil)
}
