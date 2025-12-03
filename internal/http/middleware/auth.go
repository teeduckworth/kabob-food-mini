package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Chain composes multiple Gin middleware into one.
func Chain(mws ...gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, mw := range mws {
			mw(c)
			if c.IsAborted() {
				return
			}
		}
	}
}

const userIDContextKey = "user_id"

// JWTAuth returns Gin middleware verifying Authorization Bearer tokens.
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearer(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		parsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		userID, err := extractUserID(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token subject"})
			return
		}

		c.Set(userIDContextKey, userID)
		c.Next()
	}
}

// UserIDFromContext fetches authenticated user id from Gin context.
func UserIDFromContext(c *gin.Context) (int64, bool) {
	val, ok := c.Get(userIDContextKey)
	if !ok {
		return 0, false
	}
	id, ok := val.(int64)
	return id, ok
}

func extractBearer(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.Fields(header)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}

func extractUserID(claims jwt.MapClaims) (int64, error) {
	val, ok := claims["sub"]
	if !ok {
		return 0, errors.New("sub missing")
	}
	switch v := val.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case json.Number:
		parsed, err := v.Int64()
		if err != nil {
			return 0, err
		}
		return parsed, nil
	default:
		return 0, errors.New("unexpected sub type")
	}
}
