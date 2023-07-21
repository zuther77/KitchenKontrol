package routes

import (
	controller "golang-KitchenKontrol/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/food", controller.GetAllFood())
	incomingRoutes.GET("/food/:food_id", controller.GetFood())
	incomingRoutes.POST("/food", controller.CreateFood())
	incomingRoutes.PATCH("/food/:food_id", controller.UpdateFood())
}
