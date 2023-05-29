package routes

import(
	"github.com/gin-gonic/gin"
	"golang-KitchenKontrol/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controller.GetAllUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.POST("/users/signup", controller.SignupUser())
	incomingRoutes.POST("/users/login", controller.LoginUser())
	incomingRoutes.DELETE("/users/:user_id", controller.DeleteUser())
}