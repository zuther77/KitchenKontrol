package routes

import (
	controller "golang-KitchenKontrol/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orderItem", controller.GetAllOrderItem())
	incomingRoutes.GET("/orderItem/:orderItem_id", controller.GetOrderItem())
	incomingRoutes.POST("/orderItem", controller.CreateOrderItem())
	incomingRoutes.PATCH("/orderItem/:orderItem_id", controller.UpdateOrderItem())
	incomingRoutes.GET("/orderItem/order/:order_id", controller.GetOrderItemByOrder())
}
