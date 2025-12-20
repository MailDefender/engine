package handlers

type apiError struct {
	Error string `json:"error"`
}

func newError(msg string) apiError {
	return apiError{Error: msg}
}
