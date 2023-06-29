package api

import (
	"errors"
	"objects"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// creates the JWT (JSON Web Token) and returns it
func CreateJWT(jwtSecret []byte, expiredIn int, user objects.User) (string, error) {
	currentTime := time.Now().UTC()

	expiryDate := convertToSeconds(expiredIn)
	if expiredIn == 0 {
		expiryDate = 3600 * time.Second
	}

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(expiryDate)),
		Subject:   strconv.Itoa(user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return ss, nil
}

// extracts the JWT from the auth string
func ExtractJWT(authString string) string {
	return authString[7:]
}

// checks if the JWT token is valid
func ValidateJWT(jwtToken string) (int, bool, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(LoadEnv()), nil
	})
	if err != nil {
		return 0, false, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		//check expiry date of the token
		if isExpired(claims.ExpiresAt) {
			return 0, false, errors.New("expired token")
		}
		id, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return 0, true, errors.New("error converting string to int")
		}
		return id, true, nil

	}
	return 0, true, err
}

//helper methods

// converts the expiredIn to the time in seconds
func convertToSeconds(expiredIn int) time.Duration {
	return (time.Second * time.Duration(expiredIn))
}

// checks if the token expired
func isExpired(tokenDate *jwt.NumericDate) bool {
	return time.Now().UTC().Unix() > tokenDate.Unix()
}
