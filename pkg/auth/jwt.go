package auth

import (
	"errors"
	"fmt"
	"gochatbackend/model"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtKey = []byte("FDr1VjVQiSiybYJrQZNt8Vfd7bFEsKP6vNX1brOSiWl0mAIVCxJiR4/T3zpAlBKc2/9Lw2ac4IwMElGZkssfj3dqwa7CQC7IIB+nVxiM1c9yfowAZw4WQJ86RCUTXaXvRX8JoNYlgXcRrK3BK0E/fKCOY1+izInW3abf0jEeN40HJLkXG6MZnYdhzLnPgLL/TnIFTTAbbItxqWBtkz6FkZTG+dkDSXN7xNUxlg==")

type authClaims struct {
	jwt.StandardClaims
	UserID string `json:"userId"`
}

func GenerateToken(user model.User) (string, error) {
	expiresAt := time.Now().Add(240 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, authClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: expiresAt,
		},
		UserID: user.ID.Hex(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	var claims authClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	id := claims.UserID
	return id, nil
}

func GetToken(c *gin.Context) (string, bool) {
	authValue := c.GetHeader("Authorization")
	arr := strings.Split(authValue, " ")
	if len(arr) != 2 {
		return "", false
	}
	authType := strings.Trim(arr[0], "\n\r\t")
	if strings.ToLower(authType) != "bearer" {
		return "", false
	}
	return strings.Trim(arr[1], "\n\t\r"), true
}

func ParseIdFromCtx(c *gin.Context) (primitive.ObjectID, error) {
	tokenStr, _ := GetToken(c)
	token, err := jwt.Parse(tokenStr, nil)
	if token == nil {
		fmt.Println(err)
		return primitive.NilObjectID, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := claims["userId"]
	if userId == nil {
		fmt.Println(err)
		return primitive.NilObjectID, fmt.Errorf("error when get user id")
	}
	userObjectId, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		fmt.Println(err)
		return primitive.NilObjectID, fmt.Errorf("error when parse user id")
	}
	return userObjectId, nil
}
