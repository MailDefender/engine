package engine

import (
	"fmt"

	"maildefender/engine/internal/configuration"
)

func GenerateValidationUri(token string) string {
	return fmt.Sprintf("%s/validate/%s", configuration.ValidatorPublicBaseEndpoint(), token)
}
