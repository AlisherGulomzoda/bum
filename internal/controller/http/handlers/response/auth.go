package response

// UserToken represent user token response.
type UserToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// NewUserToken creates a new user token claims.
func NewUserToken(accessToken, refreshToken string) UserToken {
	return UserToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
