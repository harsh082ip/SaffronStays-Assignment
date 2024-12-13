package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/SaffronStays-Assignment/controllers"
)

func AppRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/:room_id", controllers.MetricsControllers)
}
