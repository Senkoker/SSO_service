package jwt_token

import (
	"GRPC_Service_sso/internal/module"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewJWT(user module.User, app module.AppID, duration time.Duration) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["user.id"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app.id"] = app.Id
	token, err := jwtToken.SignedString([]byte(app.Secret))
	if err != nil {
		return "", fmt.Errorf("Problem to create JWT token")
	}
	return token, nil
}
