package dto

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type GenerateTokensRequest struct {
	GUID string `json:"guid" binding:"required,uuid"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	AccessToken  string `json:"access_token" binding:"required"`
}
type RefreshTokenData struct {
	UserGUID  string
	UserAgent string
	IP        string
}
