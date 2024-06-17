package middleware

import (
	"net/http"

	helper "github.com/ahmedgulabkhan/ChronicAid/helpers"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken, _ := c.Cookie("session_token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The user is not logged in"})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("user_name", claims.User_name)

		c.Next()
	}
}
