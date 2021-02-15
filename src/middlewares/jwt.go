package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/weeber-id/desatanjungbunga-backend/src/models"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
)

// AdminAuthorization using JWT
func AdminAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtConfig := variables.JWTConfig

		token, err := c.Cookie(jwtConfig.TokenName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "You must login before"})
			return
		}

		// decode JWT
		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(variables.JWTConfig.Key), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			return
		}

		c.Set("JWT_ROLE", int(claims["role"].(float64)))
		c.Set("JWT_ID", claims["id"])
	}
}

// GetClaims from gin context
func GetClaims(c *gin.Context) *models.JwtClaims {
	claims := new(models.JwtClaims)

	claims.Role = c.GetInt("JWT_ROLE")
	claims.ID = c.GetString("JWT_ID")

	return claims
}

// WriteAccessToken2Cookie middleware
func WriteAccessToken2Cookie(c *gin.Context, adminID string, role int) {
	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["role"] = role
	claims["id"] = adminID

	jwtConfig := variables.JWTConfig
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := at.SignedString([]byte(jwtConfig.Key))
	if err != nil {
		log.Panicf("error generate access token: %v", err)
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     jwtConfig.TokenName,
		Value:    token,
		SameSite: 4, // None
		MaxAge:   3600 * 24,
		HttpOnly: true,
		Path:     jwtConfig.Path,
		Domain:   jwtConfig.Domain,
		Secure:   jwtConfig.HTTPS,
	})
}

// DeleteAccessToken2Cookie middleware
func DeleteAccessToken2Cookie(c *gin.Context) {
	jwtConfig := variables.JWTConfig

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     jwtConfig.TokenName,
		Value:    "",
		SameSite: 4, // None
		MaxAge:   0,
		HttpOnly: true,
		Path:     jwtConfig.Path,
		Domain:   jwtConfig.Domain,
		Secure:   jwtConfig.HTTPS,
	})
}
