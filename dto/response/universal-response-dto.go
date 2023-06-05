package response

type UniversalResponse struct {
	HttpCode        int         `json:"-"`
	ResponseMessage string      `json:"responseMessage"`
	ResponseCode    string      `json:"responseCode"`
	Data            interface{} `json:"data"`
}

type BadRequestResponse struct {
	HttpCode        int         `json:"-"`
	ResponseMessage string      `json:"responseMessage"`
	ResponseCode    string      `json:"responseCode"`
	Errors          interface{} `json:"errors"`
}
