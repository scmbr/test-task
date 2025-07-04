package dto

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input data"`
}
type MessageResponse struct {
	Message string `json:"message" example:"succesfully logged out"`
}
