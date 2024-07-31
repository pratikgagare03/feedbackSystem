package helper

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        uint
	User_type  string
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, userType string, userID uint) (signedToken string, claims *SignedDetails, error error) {
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

	// refreshClaims := &SignedDetails{
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
	// 	},
	// }

	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	// signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return signedToken, claims, err
}
func UpdateAllTokens(token string, userId uint) {
	var updatedUser models.User
	updatedUser.ID = userId
	repository.Db.Find(&updatedUser)
	// updatedUser.Refresh_token = refreshToken
	repository.Db.Save(updatedUser)
	return
}

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

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "token invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token expired"
		return
	}
	return claims, msg
}
