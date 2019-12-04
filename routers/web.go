package routers

import (
	"gin/app/controller"
	"github.com/gin-gonic/gin"
)

func StuffHttp(router *gin.Engine) {
	router.Any("/", controller.Index)

	//home := router.Group("/home")
	//{
	//
	//}
}