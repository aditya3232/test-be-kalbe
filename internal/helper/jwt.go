package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var JWT_SECRET = viper.GetString("jwt.secret")

func GenerateToken(employeeId int, expiresIn int) (string, error) {
	claims := jwt.MapClaims{
		"employee": map[string]interface{}{
			"employee_id": employeeId,
		},
	}

	// expiresIn is optional, if not set, it will be 1 day
	// if expiresIn is set, it will be in days
	if expiresIn != 0 {
		claims["exp"] = time.Now().Add(time.Hour * 24 * time.Duration(expiresIn)).Unix()
	} else {
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JWT_SECRET)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
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
