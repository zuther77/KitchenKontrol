package routes

import (
	controller "golang-KitchenKontrol/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/order", controller.GetAllOrder())
	incomingRoutes.GET("/order/:order_id", controller.GetOrder())
	incomingRoutes.POST("/order", controller.CreateOrder())
	incomingRoutes.PATCH("/order/:order_id", controller.UpdateOrder())
}
