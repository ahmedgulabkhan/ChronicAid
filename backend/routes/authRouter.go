package routes

import (
	controller "github.com/ahmedgulabkhan/ChronicAid/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/api/auth/verify", controller.VerifyAuth())
	incomingRoutes.POST("/api/auth/signup", controller.SignUp())
	incomingRoutes.POST("/api/auth/login", controller.Login())
	incomingRoutes.POST("/api/auth/logout", controller.Logout())
}
