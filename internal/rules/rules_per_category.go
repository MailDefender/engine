package rules

import "maildefender/engine/internal/models"

type RulesPerCatergory struct {
	Rules map[string][]models.Rule `json:"rules"`
}

func (rpc *RulesPerCatergory) Align() []models.Rule {
	var out []models.Rule
	for c, r := range rpc.Rules {
		for _, rule := range r {
			rule.Category = c
			out = append(out, rule)
		}
	}
	return out
}
