package responses

type ResponseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ErrorResponse(message string) ResponseBody {
	return ResponseBody{Code: 1, Message: message, Data: nil}
}

func SuccessResponse(data interface{}) ResponseBody {
	return ResponseBody{Code: 0, Message: "success", Data: data}
}

func FormatResponse(code int, message string, data interface{}) ResponseBody {
	return ResponseBody{Code: code, Message: message, Data: data}
}
