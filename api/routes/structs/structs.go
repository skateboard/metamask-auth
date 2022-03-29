package structs

type MessageStruct struct {
	Message string `json:"message"`
}

type SuccessfulResponse struct {
	Data interface{} `json:"data"`
}