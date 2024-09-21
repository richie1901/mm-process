package models

type Response struct{
	ResponseCode string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Meta string `json:"meta"`
}

type CreateUserResponseSuccess struct{
	ResponseCode string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Meta User `json:"meta"`
}

type CreateUserResponseFailure struct{
	ResponseCode string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Meta User `json:"meta"`
}

type GetUsersResponseSuccess struct{
	ResponseCode string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Meta []User `json:"meta"`
}

type GetUsersResponseFailure struct{
	ResponseCode string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Meta string `json:"meta"`
}