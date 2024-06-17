package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/ahmedgulabkhan/ChronicAid/database"
	"github.com/ahmedgulabkhan/ChronicAid/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userName := c.Param("user_name")
		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_name": userName}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
