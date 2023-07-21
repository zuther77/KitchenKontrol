package main

import (
	"golang-KitchenKontrol/database"
	middleware "golang-KitchenKontrol/middleware"
	routes "golang-KitchenKontrol/routes"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	// initialize routes
	routes.FoodRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":" + port)

}
