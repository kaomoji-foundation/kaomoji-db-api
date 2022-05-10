package stdMessages

type errorBasic struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ErrorDefault(message string, data interface{}) errorBasic {
	return errorBasic{
		Status:  "error",
		Message: message,
		Data:    data,
	}
}
