package responseDTO

type Response struct {
	HttpCode         int         `json:"httpCode"`
	ResponseCode     string      `json:"responseCode"`
	ResponseDesc     string      `json:"responseDesc"`
	ResponseData     interface{} `json:"responseData"`
	MetadataResponse interface{} `json:"metadata"`
	Summary          interface{} `json:"summary"`
}

type MetadataResponse struct {
	ListUserDtoRes ListUserDtoRes `json:"listUserDtoRes"`
}

type ListUserDtoRes struct {
	Currentpage  int `json:"currentPage"`
	TotalData    int `json:"totalData"`
	TotalDataAll int `json:"totalDataAll"`
}

type MetadataSummeryResponse struct {
	ListTerdekat ListUserDtoRes `json:"listTerdekat"`
	List360      ListUserDtoRes `json:"list360"`
	ListByCity   ListUserDtoRes `json:"listByCity"`
}
