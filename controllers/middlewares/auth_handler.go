package middlewares

import (
	"net/http"
	"strings"

	controllers "github.com/ankur12345678/uptime-monitor/Controllers"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

func HandleAuth(c *gin.Context) {
	//remove loading of configs from here
	env := controllers.Ctrl.Config

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error_message": "No auth token found",
		})
		return
	}

	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) != 2 {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error_message": "No auth token found",
		})
		return
	}

	authToken := splitAuthHeader[1]
	if authToken == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error_message": "No auth token found",
		})
		return
	}

	//verify jwt
	claims := jwt.MapClaims{}
	parsedToken, _ := jwt.ParseWithClaims(authToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.JwtSecret), nil
	})

	if !parsedToken.Valid {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error_message": "Invalid Auth token",
		})
		return
	}

	//setting jit(jwt_id) in the context so that we can use it in token generation in /refresh
	emailFromClaims := claims["email"]
	email := emailFromClaims.(string)

	jtiFromClaims := claims["jti"]
	jti := jtiFromClaims.(string)

	//check in redis for this jti. if the same access token exist corresponding to this jti then return blacklisted!
	val, err := controllers.Ctrl.RedisClient.Get(controllers.Ctrl.RedisClient.Context(), jti).Result()
	if err != nil && err != redis.Nil {
		logger.Error("error getting key from redis | err: ", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error_message": "Internal server error",
		})
		return
	}
	if val == authToken {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error_message": "Please login again",
		})
		return
	}
	c.Set("accessToken", authToken)
	c.Set("email", email)
	c.Set("jti", jti)
	c.Next()

}
