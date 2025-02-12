package guard

import (
	"example/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

const MAX_DURATION = 10.0

func HasTimeElapsed(user model.User) bool {

	lastVerified := user.LastVerified
	fmt.Println("lastVerified", lastVerified)
	now := time.Now()
	fmt.Println("now", now)
	timeOfVerified := now.Sub(lastVerified).Minutes()
	fmt.Println("timeOfVerified", timeOfVerified)

	if timeOfVerified > 10.00 {
		return true
	}
	return false
}

func ValidateUserWithPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userId := ExtractToken(tokenString)
		user := model.User{UserId: userId, Token: tokenString}
		db.Find(&user)
		if HasTimeElapsed(user) {
			c.JSON(401, gin.H{"error": "Expired. Please reauthenticate."})
			c.Abort()
			return
		}
		c.Set("userId", user.UserId)
		c.Next()
	}
}
