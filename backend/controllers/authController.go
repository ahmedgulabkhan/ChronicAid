package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	helper "github.com/ahmedgulabkhan/ChronicAid/helpers"
	"github.com/ahmedgulabkhan/ChronicAid/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func VerifyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken, _ := c.Cookie("session_token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The user is not logged in"})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		var user models.User
		user.First_name = &claims.First_name
		user.User_name = claims.User_name

		c.JSON(http.StatusOK, user)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.StructPartial(user, "Email", "Password")
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login or passowrd is incorrect"})
			return
		}

		passwordIsValid, msg := verifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		signedToken, err := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		foundUser.Token = &signedToken

		cookie := &http.Cookie{
			Name:     "session_token",
			Value:    signedToken,
			MaxAge:   86400,
			Path:     "/",
			Domain:   "",
			HttpOnly: true,
			Secure:   false,
		}

		http.SetCookie(c.Writer, cookie)

		c.JSON(http.StatusOK, foundUser)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for the email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This email already exists"})
			return
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"email": user.User_name})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for the username"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This username already exists"})
			return
		}

		password := hashPassword(*user.Password)
		user.Password = &password
		user.ID = primitive.NewObjectID()
		token, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *&user.User_name)
		user.Token = &token

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		cookie := &http.Cookie{
			Name:     "session_token",
			Value:    token,
			MaxAge:   86400,
			Path:     "/",
			Domain:   "",
			HttpOnly: true,
			Secure:   false,
		}

		http.SetCookie(c.Writer, cookie)

		c.JSON(http.StatusOK, user)
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully!"})
	}
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func verifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Login or passowrd is incorrect")
		check = false
	}

	return check, msg
}
