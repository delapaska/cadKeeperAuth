package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/delapaska/cadKeeperAuth/configs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(configs.Envs.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa((userID)),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func WithJWTAuth(handlerFunc func(*gin.Context), secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := getTokenFromRequest(c)
		if tokenString == "" {
			log.Println("token is missing")
			permissionDenied(c)
			return
		}

		token, err := validateJWT(tokenString, secret)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(c)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(c)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Println("invalid token claims")
			permissionDenied(c)
			return
		}

		str, ok := claims["userID"].(string)
		if !ok {
			log.Println("userID not found in token claims")
			permissionDenied(c)
			return
		}

		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			permissionDenied(c)
			return
		}

		// Add the userID to the context
		c.Set(string(UserKey), userID)

		// Call the handler function with userID
		handlerFunc(c)
	}
}

func validateJWT(tokenString string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Use the provided secret key
		return []byte(secret), nil
	})
}

func permissionDenied(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
	c.Abort()
}

func getTokenFromRequest(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
