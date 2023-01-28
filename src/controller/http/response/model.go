package response

type Map map[string]interface{}

// Error error
// swagger:model Error
type AppError struct {
	code             int
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Params           Map    `json:"params,omitempty"`
}

type SuccessResponse struct {
	code    int
	Message string `json:"message"`
	Params  Map    `json:"params,omitempty"`
}
