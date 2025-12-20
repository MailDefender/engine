package notification

import (
	"testing"

	"maildefender/engine/internal/constants"
	"maildefender/engine/internal/models"
)

var wrapperTests = []struct {
	msg                messageWrapper
	expectedBadgeType  string
	expectedBadgeLabel string
}{
	{
		msg: messageWrapper{
			Entry: models.MessageHistory{
				Source:      string(constants.DefaultMailbox),
				Destination: "MAILBOX.RANDOM",
			},
		},
		expectedBadgeType:  "new",
		expectedBadgeLabel: "New",
	},
	{
		msg: messageWrapper{
			Entry: models.MessageHistory{
				Source:      string(constants.DefaultMailbox),
				Destination: string(constants.BlackListedMailbox),
			},
		},
		expectedBadgeType:  "blacklisted",
		expectedBadgeLabel: "Blacklisted",
	},
	{
		msg: messageWrapper{
			Entry: models.MessageHistory{
				Source:      string(constants.DefaultMailbox),
				Destination: string(constants.PendingMailbox),
			},
		},
		expectedBadgeType:  "pending",
		expectedBadgeLabel: "Pending",
	},
	{
		msg: messageWrapper{
			Entry: models.MessageHistory{
				Source:      string(constants.PendingMailbox),
				Destination: "MAILBOX.RANDOM",
			},
		},
		expectedBadgeType:  "resolved",
		expectedBadgeLabel: "Resolved",
	},
}

func TestMessageWrapper(t *testing.T) {
	for _, test := range wrapperTests {
		badgeType, badgeLabel := test.msg.badgeType(), test.msg.badgeLabel()

		if badgeLabel != test.expectedBadgeLabel {
			t.Errorf("Wrong badge label, expected '%s', got '%s'", test.expectedBadgeLabel, badgeLabel)
		}

		if badgeType != test.expectedBadgeType {
			t.Errorf("Wrong badge type, expected '%s', got '%s'", test.expectedBadgeType, badgeType)
		}
	}
}
