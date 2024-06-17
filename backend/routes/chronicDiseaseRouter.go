package routes

import (
	controller "github.com/ahmedgulabkhan/ChronicAid/controllers"
	"github.com/ahmedgulabkhan/ChronicAid/middleware"

	"github.com/gin-gonic/gin"
)

func ChronicDiseaseRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/api/diseases/symptoms", controller.CheckSymptomsForChronicDisease())

	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/api/diseases/:user_name", controller.GetUserChronicDiseases())
	incomingRoutes.GET("/api/diseases/metrics/all/:disease_name", controller.GetMetricNamesForDisease())
	incomingRoutes.GET("/api/diseases/metrics/alllast7days/:user_name", controller.GetAllUserMetricsForLast7Days())
	incomingRoutes.GET("/api/diseases/metrics/:user_name/:disease_name/today", controller.GetUserChronicDiseaseMetricsForToday())
	incomingRoutes.GET("/api/diseases/metrics/:user_name/:disease_name/last7days", controller.GetUserChronicDiseaseMetricsForLast7Days())
	incomingRoutes.GET("/api/diseases/metrics/:user_name/:disease_name/last28days", controller.GetUserChronicDiseaseMetricsForLast28Days())

	incomingRoutes.POST("/api/diseases/:user_name", controller.AddNewUserChronicDisease())
	incomingRoutes.POST("/api/diseases/metrics/:user_name/:disease_name", controller.AddNewChronicDiseaseMetric())
}
