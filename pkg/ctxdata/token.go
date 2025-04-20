package ctxdata

import "github.com/golang-jwt/jwt"

const Identiy = "myIm.com"

func GetJwtToken(secret string, iat, seconds int64, uid string) (string, error) {
	claims := jwt.MapClaims{}
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[Identiy] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secret))
}
