package model

type DefaultResponse struct {
	Message string `json:"message"`
}

func Response(message string) DefaultResponse {
	return DefaultResponse{
		Message: message,
	}
}
