package routers

import (
	"gin/app/controller"
	"github.com/gin-gonic/gin"
)

func StuffApi(router *gin.Engine) {
	router.Any("/api", controller.Index)

	api := router.Group("/api")
	{
		api.GET("/person", controller.ListPerson)
		api.GET("/person/:id", controller.FindPerson)
		api.PUT("/person/:id", controller.EditPerson)
		api.DELETE("/person/:id", controller.DelPerson)
		api.POST("/person", controller.AddPerson)
	}
}