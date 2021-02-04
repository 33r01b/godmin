package dto

type Token struct {
	AccessToken         string
	RefreshToken        string
	AccessUuid          string
	RefreshUuid         string
	AccessTokenExpires  int64
	RefreshTokenExpires int64
}
