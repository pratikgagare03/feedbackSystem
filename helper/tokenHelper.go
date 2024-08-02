package helper

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pratikgagare03/feedback/logger"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        uint
	User_type  string
	jwt.StandardClaims
}

// SECRET_KEY is the key used to sign the JWT token.
var SECRET_KEY = os.Getenv("SECRET_KEY")

// GenerateAccessToken generates a JWT token with the given email, first name, last name, user type, and user ID.
func GenerateAccessToken(email, firstName, lastName, userType string, userID uint) (signedToken string, claims *SignedDetails, error error) {
	// Create a new SignedDetails struct with the given email, first name, last name, user type, and user ID.
	claims = &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_type:  userType,
		Uid:        userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	// Create a new JWT token with the SigningMethodHS256 signing method and the claims.
	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		logger.Logs.Panic().Msg("error while signing token")
	}

	// Return the signed token, the claims, and any error that occurred.
	return signedToken, claims, err
}

// validateToken validates the given signed token and returns the claims and a message.
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	// Check if the token is valid.
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "token invalid"
		return
	}
	// Check if the token is expired.
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token expired"
		return
	}
	return claims, msg
}
