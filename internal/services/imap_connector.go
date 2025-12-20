package services

import (
	"maildefender/engine/internal/configuration"
	"maildefender/engine/internal/third_party/models"
)

type FetchEmailsIn struct {
	Mailbox string
	Sender  *string
}

type fetchEmailsOut struct {
	Error    string          `json:"error"`
	Count    int             `json:"count"`
	Messages models.Messages `json:"messages"`
}

func FetchEmails(in FetchEmailsIn) (models.Messages, int, error) {
	var data fetchEmailsOut

	queryParams := map[string]string{
		"mailbox": in.Mailbox,
	}

	if in.Sender != nil {
		queryParams["sender"] = *in.Sender
	}

	httpCode, err := doReq(
		configuration.ImapConnectorBaseEndpoint(),
		request{
			reqType:     get,
			endpoint:    "/v1/imap-connector/emails",
			queryParams: queryParams,
		},
		&data,
	)
	return data.Messages, httpCode, err
}

type createMailboxIn struct {
	Mailbox string `json:"mailbox"`
}

type createMailboxOut struct {
	Error string `json:"error"`
}

func CreateMailbox(mailbox string) (createMailboxOut, int, error) {
	var data createMailboxOut
	httpCode, err := doReq(
		configuration.ImapConnectorBaseEndpoint(),
		request{
			reqType:  post,
			endpoint: "/v1/imap-connector/mailboxes",
			body: createMailboxIn{
				Mailbox: mailbox,
			},
		},
		&data,
	)
	return data, httpCode, err
}

type MoveEmailIn struct {
	MessageID     string `json:"messageId"`
	Source        string `json:"source"`
	Destination   string `json:"destination"`
	CreateMailbox bool   `json:"createMailbox"`
}

type MoveEmailOut struct {
	Error string `json:"error"`
}

func MoveEmail(in MoveEmailIn) (MoveEmailOut, int, error) {
	var data MoveEmailOut
	httpCode, err := doReq(
		configuration.ImapConnectorBaseEndpoint(),
		request{
			reqType:  post,
			endpoint: "/v1/imap-connector/emails/move",
			body:     in,
		},
		&data,
	)

	return data, httpCode, err
}
