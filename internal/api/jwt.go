package api

import (
	"db"
	"errors"
	"objects"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// creates the JWT Access Token (JSON Web Token) and returns it
func CreateJWTAccess(jwtSecret []byte, user objects.User) (string, error) {
	currentTime := time.Now().UTC()
	expiryDate := 3600 * time.Second

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
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

// creates the JWT Refresh (JSON Web Token) and returns it
func CreateJWTRefresh(jwtSecret []byte, user objects.User) (string, error) {
	currentTime := time.Now().UTC()
	expiryDate := time.Hour * 24

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
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

// checks if the JWT Access token is valid
func ValidateJWTAccess(jwtToken string) (int, bool, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(LoadEnv()), nil
	})
	if err != nil {
		return 0, false, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if isExpired(claims.ExpiresAt) {
			return 0, false, errors.New("expired token")
		}
		if isRefreshToken(claims) {
			return 0, false, errors.New("cannot use refresh token as access token")
		}
		id, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return 0, true, errors.New("error converting string to int")
		}
		return id, true, nil

	}
	return 0, true, err
}

// checks if the JWT Refresh token is valid
func ValidateJWTRefresh(jwtToken string, db *db.DB, database objects.DBStructure) (objects.User, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(LoadEnv()), nil
	})
	if err != nil {
		return objects.User{}, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if isExpired(claims.ExpiresAt) {
			return objects.User{}, errors.New("expired token")
		}
		if !isRefreshToken(claims) {
			return objects.User{}, errors.New("invalid token")
		}
		if db.IsTokenRevoked(jwtToken, database) {
			return objects.User{}, errors.New("token has been revoked")
		}
		id, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return objects.User{}, errors.New("error converting string to int")
		}
		user, err := db.GetUserByID(database, id)
		if err != nil {
			return objects.User{}, err
		}
		return user, nil
	}
	return objects.User{}, errors.New("invalid token or nil claims")
}

// helper methods

// converts the expiredIn to the time in seconds
func convertToSeconds(expiredIn int) time.Duration {
	return (time.Second * time.Duration(expiredIn))
}

// checks if the token expired
func isExpired(tokenDate *jwt.NumericDate) bool {
	return time.Now().UTC().Unix() > tokenDate.Unix()
}

// extracts the JWT from the auth string
func ExtractJWT(authString string) string {
	return authString[7:]
}

// determines if the refreshToken is present
func isRefreshToken(claims *jwt.RegisteredClaims) bool {
	return claims.Issuer == "chirpy-refresh"
}
