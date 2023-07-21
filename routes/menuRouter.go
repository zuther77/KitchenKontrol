package routes

import (
	controller "golang-KitchenKontrol/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/menu", controller.GetAllMenu())
	incomingRoutes.GET("/menu/:menu_id", controller.GetMenu())
	incomingRoutes.POST("/menu", controller.CreateMenu())
	incomingRoutes.PATCH("/menu/:menu_id", controller.UpdateMenu())
}
