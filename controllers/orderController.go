package controller

import (
	"context"
	"fmt"
	"golang-KitchenKontrol/database"
	"golang-KitchenKontrol/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)

func GetAllOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching orders"})
		}

		var allOrders []bson.M

		if err = result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders)

	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching order"})
		}
		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var table models.Table
		var order models.Order

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationError := validate.Struct(order)

		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		if order.Table_id != nil {

			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel()

			if err != nil {
				msg := fmt.Sprintf("Table not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			order.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			order.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

			order.ID = primitive.NewObjectID()

			order.Order_id = order.ID.Hex()

			result, insertErr := orderCollection.InsertOne(ctx, order)

			if insertErr != nil {
				msg := fmt.Sprintf("Order item was not created")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}

			defer cancel()

			c.JSON(http.StatusOK, result)
		}

	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
		var order models.Order
		var table models.Table

		orderId := c.Param("order_id")

		var updateObj primitive.D

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if order.Table_id != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel()

			if err != nil {
				msg := fmt.Sprintf("Table not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			}

			updateObj = append(updateObj, bson.E{"table_id", order.Table_id})
		}

		order.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", order.Updated_At})

		upsert := true

		filter := bson.M{"order_id": orderId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{"$set", updateObj}},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprintf("Order update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}

func OrderItemOrderCreator(order models.Order) string {

	order.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()
	orderCollection.InsertOne(ctx, order)
	defer cancel()

	return order.Order_id
}
