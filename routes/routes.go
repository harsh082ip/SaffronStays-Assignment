package routes

import "github.com/gin-gonic/gin"

func AppRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/:room_id")
}
