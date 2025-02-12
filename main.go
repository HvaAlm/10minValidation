package main

import (
	"example/database"
	"example/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)
import "example/guard"

func someEndPoint(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hi, you can really enter :)"})
}

type passValidator struct {
	Password string `json:"password" binding:"required,min=4"`
	UserId   string `json:"userId" binding:"required,min=4"`
}

func userAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body passValidator
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		userId := body.UserId
		user := model.User{UserId: userId}
		db.First(&user)
		if guard.ValidateUserWithPassword(user.Password, body.Password) {
			user.LastVerified = time.Now()
			db.Where("user_id = ?", userId).Save(&user)
			c.JSON(http.StatusOK, gin.H{
				"message": "User registered successfully!",
				"user":    user,
			})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"err": "Forbidden"})
		}
	}
}
func createUser(db *gorm.DB) {
	pass, _ := guard.HashPassword("qazwssxceddcc")
	db.Create(model.User{
		LastVerified: time.Now(),
		UserId:       "q2w3e4r",
		Password:     pass,
		Token:        guard.CreateToken("q2w3e4r"),
	})
}
func main() {
	r := gin.Default()
	db, err := database.DBConnect()
	if err != nil {
		log.Fatal("error: " + err.Error())
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal(err.Error())
	}

	r.GET("/welcome", guard.AuthMiddleware(db), someEndPoint)

	r.POST("/password/validation", userAuth(db))

	err = r.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
