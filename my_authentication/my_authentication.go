package my_authentication

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var secretKey = []byte("your-secret-key")

var userClaims = UserClaims{
	UserID:   1,
	UserName: "User1",
	Role:     "Admin",
	RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "my-app",
		Subject:   "user-token",
		ID:        "1",
	},
}

func getJWTToken() *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
}

func GETJWTTokenString() (string, error) {
	token := getJWTToken()
	return token.SignedString(secretKey)
}

func getTokenFromHeader(c *gin.Context) (string, error) {
	authString := c.GetHeader("Authorization")
	if authString == "" {
		return "", errors.New("authen failed: Not found any authorization header from header")
	}

	if !strings.HasPrefix(authString, "Bearer ") {
		return "", errors.New("invalid Auth Type")
	}

	token := strings.TrimSpace(authString[len("Bearer "):])
	return token, nil
}

func verifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra xem thuật toán ký có phải là HS256 không
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Kiểm tra xem token có hợp lệ không
	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := getTokenFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, err := verifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Lưu thông tin người dùng vào context để các handler có thể sử dụng
		c.Set("user", claims)
		c.Next()
	}
}
