package routes

import (
	controller "github.com/ahmedgulabkhan/ChronicAid/controllers"
	"github.com/ahmedgulabkhan/ChronicAid/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/api/users/:user_name", controller.GetUser())
}
