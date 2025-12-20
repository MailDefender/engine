package models

import "time"

type Message struct {
	MessageID string    `json:"messageId"`
	Headers   Headers   `json:"headers"`
	From      Contacts  `json:"from,omitempty"`
	To        Contacts  `json:"to,omitempty"`
	Cc        Contacts  `json:"cc,omitempty"`
	Subject   string    `json:"subject,omitempty"`
	Date      time.Time `json:"date"`
}

type Messages []Message
