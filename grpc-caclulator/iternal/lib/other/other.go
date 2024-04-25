package other

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func RemoveEmptyStrings(s []string) []string {
	result := make([]string, 0)
	for _, str := range s {
		if str != " " {
			result = append(result, str)
		}
	}
	return result
}
func ParseToken(tokenString string) (float64, time.Time, error) {
	// Распарсивание токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("culc"), nil
	})
	if err != nil {
		return 0, time.Time{}, err
	}

	// Проверка на валидность токена
	if !token.Valid {
		return 0, time.Time{}, errors.New("invalid token")
	}

	// Извлечение данных из полезной нагрузки (claims)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, time.Time{}, errors.New("invalid claims")
	}

	// Извлечение uid из полезной нагрузки
	uid, ok := claims["uid"].(float64)
	fmt.Println(uid)
	if !ok {

		fmt.Println(uid)
		return 0, time.Time{}, errors.New("invalid uid")
	}

	// Извлечение времени истечения (exp) из полезной нагрузки и преобразование во время
	expInt, ok := claims["exp"].(float64)

	if !ok {
		return 0, time.Time{}, errors.New("invalid exp")
	}
	//timelive, ok := claims["timelive"].(float64)

	//if !ok {
	//	return 0, time.Time{}, time.Time{}, errors.New("invalid timeliive")
	//}
	exp := time.Unix(int64(expInt), 0)
	//timel := time.Unix(int64(timelive), 0)

	return uid, exp, nil
}
