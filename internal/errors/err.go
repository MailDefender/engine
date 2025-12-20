package errors

import "errors"

var (
	ErrCategoryRuleEmptyName        = errors.New("empty category rule name")
	ErrCategoryRuleEmptyDestination = errors.New("empty category rule destination")

	ErrCriterionInvalidType  = errors.New("invalid criterion type")
	ErrCriterionInvalidInput = errors.New("invalid criterion input")
	ErrCriterionEmptyValues  = errors.New("empty criterion values")

	ErrCannotFindMatchingRules = errors.New("cannit find matching rules")
	ErrNoRuleFound             = errors.New("no rule found")
	ErrTooManyRulesFound       = errors.New("too many rules found")

	ErrCannotSaveHistory = errors.New("cannot save history")

	ErrExpiredToken          error = errors.New("expired token")
	ErrAlreadyValidatedToken error = errors.New("token already validated")
)
