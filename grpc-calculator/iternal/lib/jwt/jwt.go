package jwtoken

import (
	"culc/iternal/model"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func NewToken(user model.User, tokentl time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["timelive"] = tokentl
	claims["exp"] = time.Now().Add(tokentl).Unix()

	tokenString, err := token.SignedString([]byte("culc"))
	if err != nil {
		return "", fmt.Errorf("error token made")
	}

	return tokenString, nil
}
