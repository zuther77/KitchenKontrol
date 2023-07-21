package controller

import (
	"context"
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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetAllTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := tableCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching orders"})
			return
		}

		var allTable []bson.M

		if err = result.All(ctx, &allTable); err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, allTable)

	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
		tableId := c.Param("table_id")
		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching order"})
		}
		c.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(table)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		table.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, insertErr := tableCollection.InsertOne(ctx, table)
		defer cancel()

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating table"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)

		tableId := c.Param("table_id")

		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if table.Table_number != nil {
			updateObj = append(updateObj, bson.E{"table_number", *table.Table_number})
		}

		if table.Number_of_guest != nil {
			updateObj = append(updateObj, bson.E{"number_of_guest", *table.Number_of_guest})
		}

		table.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", table.Updated_At})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"table_id": tableId}

		result, err := tableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating table"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}
