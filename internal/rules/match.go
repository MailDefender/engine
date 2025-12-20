package rules

import (
	"regexp"

	"maildefender/engine/internal/errors"
	"maildefender/engine/internal/models"
	thirdModels "maildefender/engine/internal/third_party/models"
)

func FindMatch(rules models.Rules, mail thirdModels.Message) ([]models.Rule, error) {
	var output []models.Rule

	for _, rule := range rules {
		match := false
		for _, criterion := range rule.Criteria {
			m, err := isCriterionMatch(criterion, mail)
			if err != nil {
				return nil, err
			}

			if m {
				match = true
				break
			}
		}
		if match {
			output = append(output, rule)
			break
		}

	}

	return output, nil
}

func isCriterionMatch(criterion models.Criterion, mail thirdModels.Message) (bool, error) {
	// First try of this criterion, we must build regexes
	if criterion.BuiltValues == nil || len(criterion.BuiltValues) == 0 {
		for _, v := range criterion.Values {
			r, err := regexp.Compile(v)
			if err != nil {
				return false, err
			}
			criterion.BuiltValues = append(criterion.BuiltValues, r)
		}
	}

	var eval evaluator

	switch criterion.Input {
	case "sender":
		eval.inputs = []string{mail.From[0].Email}
		break
	case "recipients":
		for _, rcp := range mail.To {
			eval.inputs = append(eval.inputs, rcp.Email)
		}
		break
	default:
		return false, errors.ErrCriterionInvalidInput
	}

	eval.shouldMatch = true
	eval.regex = criterion.BuiltValues
	if criterion.Count != nil {
		eval.shoudExactCount = criterion.ShouldExactCount
		eval.count = *criterion.Count
	}

	switch criterion.Type {
	case "contains":
		break
	default:
		return false, errors.ErrCriterionInvalidType
	}

	if eval.evaluate() {
		return true, nil
	}

	// TODO: manage sub-criteria

	return false, nil
}
