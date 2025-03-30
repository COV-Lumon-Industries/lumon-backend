package response

type ServerResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func NewServerResponse(data interface{}) ServerResponse {
	return ServerResponse{
		Status: "server-error",
		Data:   data,
	}
}
