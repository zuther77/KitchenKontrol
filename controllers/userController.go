package controller

import (
	"context"
	"fmt"
	"golang-KitchenKontrol/database"
	helpers "golang-KitchenKontrol/helpers"
	"golang-KitchenKontrol/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollecton *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// pagination
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		// mongodb stages
		matchStage := bson.D{{"$match", bson.D{{}}}}
		projectStage := bson.D{{
			"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}

		// aggregate data
		result, err := userCollecton.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro while listing user items "})
			return
		}

		var allUsers []bson.M

		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}

		// send response
		c.JSON(http.StatusOK, allUsers)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var userId = c.Param("userId")

		var user models.User

		err := userCollecton.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("user item not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, user)

	}
}

func SignupUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		// convert JSON data from Postman to Go struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// validate data on user struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// check if user already exists in database
		count1, err := userCollecton.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for email"})
		}

		if count1 > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user email  already exists"})
			return
		}

		// hash password
		password := HashPassword(*user.Password)
		user.Password = &password

		// check if phone number already exists in database
		count2, err := userCollecton.CountDocuments(ctx, bson.M{"phone_number": user.Phone_number})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for email"})
		}

		if count2 > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "phone number  already exists"})
			return
		}

		// meta data for user
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// generate token and referesh token
		token, refresh_token, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_token = &refresh_token

		// insert user into database
		result, insertErr := userCollecton.InsertOne(ctx, user)

		if insertErr != nil {
			msg := fmt.Sprintf("user item not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// return stauts OK
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		// convert JSON data from Postman to Go struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// check if user exists in database
		err := userCollecton.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found. Please signup"})
			return
		}

		// compare password with hashed password
		passwordMatch, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordMatch {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		// generate token and referesh token
		token, refresh_token, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

		// update token and referesh token in database
		helpers.UpdateAllTokens(token, refresh_token, foundUser.User_id)

		// return stauts OK
		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(passowrd string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passowrd), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassoword string, password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(userPassoword))
	if err != nil {
		return false, fmt.Sprintf("password does not match")
	}
	return true, fmt.Sprintf("password match")
}
