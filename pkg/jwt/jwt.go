package jwt

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// ExtractClaims extracts roles from the claims of JWT token
func ExtractClaims(tokenString string, signingKey []byte) (jwtgo.MapClaims, error) {
	claims := jwtgo.MapClaims{}
	if tokenString == "" {
		claims["role"] = "unauthorized"
		return claims, nil
	}
	if strings.Contains(tokenString, "Basic") {
		claims["role"] = "unauthorized"
		return claims, nil
	}
	token, err := jwtgo.ParseWithClaims(tokenString, claims, func(token *jwtgo.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !(ok && token.Valid) {
		err = fmt.Errorf("invalid jwt token")
		return nil, err
	}

	return claims, nil
}

// ExtractFromClaims extracts the key from jwt claim's metadata
func ExtractFromClaims(key, accessToken string, signingKey []byte) (interface{}, error) {

	claims, err := ExtractClaims(accessToken, signingKey)
	if err != nil {
		log.Println("could not extract claims:", err)
		return "", err
	}

	if _, ok := claims[key]; !ok {
		return nil, errors.New("could not find claims for key: " + key)
	}

	return claims[key], nil

}

func ExtractUserIDFromToken(c *gin.Context, secretKey []byte) (string, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return "", errors.New("authorization header is missing")
	}
	
	userId, err := ExtractFromClaims("id", token, secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to extract user ID from token: %w", err)
	}

	userIDStr, ok := userId.(string)
	if !ok {
		return "", errors.New("invalid user ID type")
	}

	return userIDStr, nil
}


// GenerateNewJWTToken generates a new JWT token
func GenerateNewJWTToken(tokenMetadata map[string]string, tokenExpireTime time.Duration, signingKey string) (string, error) {

	// Create a new claims.
	claims := jwt.MapClaims{}

	for key, value := range tokenMetadata {
		claims[key] = value
	}

	claims["iat"] = time.Now().Unix()
	claims["expires"] = time.Now().Add(tokenExpireTime).Unix()

	// Create a new JWT access token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token.
	t, err := token.SignedString([]byte(signingKey))
	if err != nil {
		// Return error, it JWT token generation failed.
		return "", err
	}

	return t, nil
}

func ParseToken(tokenString string, signingKey []byte) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return signingKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["userID"].(string)
		if !ok {
			return "", errors.New("userID not found in token")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}
