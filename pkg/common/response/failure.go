package response

type FailureResponse struct {
	Status  string `json:"status"`
	Message string `json:"data"`
}

func NewFailureResponse(message string) FailureResponse {
	return FailureResponse{
		Status:  "failure",
		Message: message,
	}
}
