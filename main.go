package main

import (
	"os"
	"golang-KitchenKontrol/database"
	"golang-KitchenKontrol/routes"
	"golang-KitchenKontrol/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
// var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
// var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
// var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")
// var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
// var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
// var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")


func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.Routes(router) 
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