package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

type jwtMiddleware struct {
	logger *zap.SugaredLogger
}
type JWTMiddleware interface {
	Authorize() gin.HandlerFunc
}

func NewJWTMiddleware() JWTMiddleware {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	return &jwtMiddleware{sugar}
}

func (j *jwtMiddleware) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		j.logger.Infof("Starting authorization middleware")
		jwtSecret := os.Getenv("JWT_SECRET")
		jwtToken, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			j.logger.Warnf("Error getting cookie token: " + err.Error())
			return
		}
		token, err := jwt.ParseWithClaims(
			jwtToken,
			&jwt.MapClaims{},
			func(token *jwt.Token) (interface{}, error) { return []byte(jwtSecret), nil })
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			j.logger.Warnf("Error parsing token")
			return
		}
		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			j.logger.Warnf("Token is invalid")
			return
		}
		expire, err := claims.GetExpirationTime()
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			j.logger.Warnf("Token is invalid")
			return
		}
		if expire.Before(time.Now()) {
			c.AbortWithStatus(http.StatusUnauthorized)
			j.logger.Warnf("Token is invalid")
			return
		}
		c.Next()
	}
}
