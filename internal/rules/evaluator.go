package rules

import "regexp"

type evaluator struct {
	inputs          []string
	regex           []*regexp.Regexp
	shouldMatch     bool
	count           int
	shoudExactCount bool
}

func (re evaluator) evaluate() bool {
	count := 0

	for _, input := range re.inputs {
		for _, reg := range re.regex {
			if reg.Match([]byte(input)) {
				count++
			}
		}
	}

	if re.shoudExactCount && re.count == count {
		return re.shouldMatch
	}

	if count > 0 {
		return re.shouldMatch
	}

	return !re.shouldMatch
}
