package configuration

import "maildefender/engine/internal/utils"

var (
	imapConnectorBaseEndpoint   string = utils.GetEnvString("IMAP_CONNECTOR_BASE_ENDPOINT", "")
	notifierBaseEndpoint        string = utils.GetEnvString("NOTIFIER_BASE_ENDPOINT", "")
	validatorPublicBaseEndpoint string = utils.GetEnvString("VALIDATOR_PUBLIC_BASE_ENDPOINT", "")
	databaseDns                 string = utils.GetEnvString("DATABASE_DNS", "")
	loopDelay                   int    = utils.GetEnvInt("LOOP_DELAYS_SECS", 5)
	rulesRefreshDelay           int    = utils.GetEnvInt("RULES_REFRESH_PERIOD_SECS", 300)
	skipReputationCheck         bool   = utils.GetEnvBool("SKIP_REPUTATION_CHECK", false)
	enableDailyMailbox          bool   = utils.GetEnvBool("ENABLE_DAILY_MAILBOX", true)
	enableDailyRecap            bool   = utils.GetEnvBool("ENABLE_DAILY_RECAP", false)
	dailyRecapRecipient         string = utils.GetEnvString("DAILY_RECAP_RECIPIENT", "")
	rulesDirectory              string = utils.GetEnvString("RULES_DIR", "")
)

func ImapConnectorBaseEndpoint() string {
	return imapConnectorBaseEndpoint
}

func NotifierBaseEndpoint() string {
	return notifierBaseEndpoint
}

func LoopDelay() int {
	return loopDelay
}

func ValidatorPublicBaseEndpoint() string {
	return validatorPublicBaseEndpoint
}

func DatabaseDNS() string {
	return databaseDns
}

func RulesRefreshDelay() int {
	return rulesRefreshDelay
}

func SkipReputationCheck() bool {
	return skipReputationCheck
}

func EnableDailyRecap() bool {
	return enableDailyRecap
}

func DailyRecapRecipient() string {
	return dailyRecapRecipient
}

func RulesDirectory() string {
	return rulesDirectory
}

func DailyMailboxEnabled() bool {
	return enableDailyMailbox
}
