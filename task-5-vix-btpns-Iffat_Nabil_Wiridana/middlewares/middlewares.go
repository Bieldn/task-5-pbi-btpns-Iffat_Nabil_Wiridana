package middlewares

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"vix-btpns/helpers"
	"vix-btpns/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := authenticateJWTToken(helpers.GetEnv("JWT_SECRET_KEY"), c)
		if err != nil {
			errorMessages := gin.H{"errors": err.Error()}
			response := helpers.APIResponse("unauthorized token", http.StatusUnauthorized, "error", errorMessages)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		encodedUserID := claim["sub"].(string)
		decodeUserID, _ := base64.StdEncoding.DecodeString(encodedUserID)
		if err != nil {
			response := helpers.APIResponse("unauthorized token", http.StatusUnauthorized, "error", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		decryptedUserID, err := helpers.Decrypt([]byte(decodeUserID))
		if err != nil {
			response := helpers.APIResponse("unauthorized token", http.StatusUnauthorized, "error", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		var user models.User
		err = db.Where("id = ?", decryptedUserID).Find(&user).Error
		if err != nil {
			response := helpers.APIResponse("unauthorized token", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}

// AuthenticateJWTToken is the main function to
// verify the JWT token from a request and it returns the claims
func authenticateJWTToken(secretKey string, c *gin.Context) (map[string]interface{}, error) {
	jwtToken, err := extractJWTToken(c)

	if err != nil {
		return nil, errors.New("Failed get JWT token")
	}

	claims, err := helpers.ParseJWT(jwtToken, secretKey)

	if err != nil {
		return nil, errors.New("Failed to parse token")
	}

	return claims, nil
}

// ExtractJWTToken extracts bearer token from Authorization header
func extractJWTToken(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		return "", errors.New("Could not find token")
	}

	tokenString, err := stripTokenPrefix(tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Strips 'Token' or 'Bearer' prefix from token string
func stripTokenPrefix(tok string) (string, error) {
	// split token to 2 parts
	tokenParts := strings.Split(tok, " ")

	if len(tokenParts) < 2 {
		return tokenParts[0], nil
	}

	return tokenParts[1], nil
}
