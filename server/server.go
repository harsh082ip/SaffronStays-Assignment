package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/SaffronStays-Assignment/routes"
)

const (
	WEBPORT = ":8000"
)

func RunServer() {
	router := gin.Default()
	routes.AppRoutes(router)

	if err := router.Run(WEBPORT); err != nil {
		log.Println("Cannot run server on:", WEBPORT, "\nerror msg:", err.Error())
	}
}
