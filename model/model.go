package model

// DefaultResponse is the default response template
type DefaultResponse struct {
	Message string `json:"message"`
}

// Response is a utility method to ease the creation of the default response
func Response(message string) DefaultResponse {
	return DefaultResponse{
		Message: message,
	}
}
