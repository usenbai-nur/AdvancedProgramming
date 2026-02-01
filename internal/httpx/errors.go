package httpx

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Err(code, message string) APIError {
	return APIError{Code: code, Message: message}
}
