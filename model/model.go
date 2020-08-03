package model

// Default response template
type DefaultResponse struct {
	Message string `json:"message"`
}

// Utility method to ease the creation of the default response
func Response(message string) DefaultResponse {
	return DefaultResponse{
		Message: message,
	}
}
