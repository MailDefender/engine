package constants

type MailboxName string

const (
	DefaultMailbox       MailboxName = "INBOX"
	MailDefenderMailbox  MailboxName = "INBOX.MailDefender"
	DailyMailbox         MailboxName = MailDefenderMailbox + "._Daily"
	PendingMailbox       MailboxName = MailDefenderMailbox + "._Pending"
	BlackListedMailbox   MailboxName = MailDefenderMailbox + "._Blacklisted"
	UncategorizedMailbox MailboxName = MailDefenderMailbox + "._Uncategorized"
)
