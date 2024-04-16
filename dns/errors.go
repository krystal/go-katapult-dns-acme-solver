package dns

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
