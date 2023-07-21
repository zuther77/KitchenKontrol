package routes

import (
	controller "golang-KitchenKontrol/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controller.GetAllUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.POST("/users/signup", controller.SignupUser())
	incomingRoutes.POST("/users/login", controller.LoginUser())
}
