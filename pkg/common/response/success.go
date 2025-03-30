package response

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Status: "success",
		Data:   data,
	}
}
