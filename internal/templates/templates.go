package templates

import (
	_ "embed"
	html "html/template"
)

var (
	//go:embed first-email-validation.html
	firstEmailValidationContent string

	//go:embed second-email-validation.html
	secondEmailValidationContent string

	//go:embed operation_recap.html
	operationRecapContent string
)

var (
	firstEmailValidationTemplate  *html.Template = html.Must(html.New("").Parse(firstEmailValidationContent))
	secondEmailValidationTemplate *html.Template = html.Must(html.New("").Parse(secondEmailValidationContent))
	operationRecapTemplate        *html.Template = html.Must(html.New("").Parse(operationRecapContent))
)

func FirstEmailValidationTemplate() *html.Template {
	return firstEmailValidationTemplate
}

func SecondEmailValidationTemplate() *html.Template {
	return secondEmailValidationTemplate
}

func OperationRecapTemplate() *html.Template {
	return operationRecapTemplate
}
