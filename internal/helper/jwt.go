package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var (
	JWT_EXPIRES_IN = 2
	JWT_SECRET     = []byte(viper.GetString("jwt.secret"))
)

func GenerateToken(employeeId int) (string, time.Time, error) {
	claims := jwt.MapClaims{
		"employee": map[string]interface{}{
			"employee_id": employeeId,
		},
	}

	// Set the expiration time for the token, in this case when JWT_EXPIRES_IN is 2, then the token will expire in 2 days
	expirationTime := time.Now().Add(time.Hour * 24 * time.Duration(int(JWT_EXPIRES_IN)))
	claims["exp"] = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JWT_SECRET)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expirationTime, nil
}

func ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return token, err
	}

	return token, nil

}

func GetEmployeeIDFromToken(encodedToken string) (int, error) {
	// Parse the JWT token
	token, err := ValidateToken(encodedToken)
	if err != nil {
		return 0, err
	}

	// Extract the user ID from the "sub" claim
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		employeeId, ok := claims["employee"].(map[string]interface{})["employee_id"].(float64)
		if !ok {
			return 0, errors.New("invalid token claims")
		}

		return int(employeeId), nil
	} else {
		return 0, errors.New("invalid token claims")
	}
}
