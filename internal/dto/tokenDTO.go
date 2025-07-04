package dto

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type GenerateTokensRequest struct {
	UserGUID string `json:"user_guid" binding:"required,uuid"`
}
type UserGUIDResponse struct {
	UserGUID string `json:"user_guid"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	AccessToken  string `json:"access_token" binding:"required"`
}
type LogoutRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}
type RefreshTokenData struct {
	UserGUID  string
	UserAgent string
	IP        string
}
