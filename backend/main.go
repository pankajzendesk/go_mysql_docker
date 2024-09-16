package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name    string `json:"name" form:"Name"`
	Surname string `json:"surname" form:"Surname"`
	Age     int    `json:"age" form:"Age"`
}

type UserWithSerial struct {
	SerialNumber int
	User
}

var db *gorm.DB

func init() {
	var err error
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ":3306)/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&User{})
}

func main() {
	r := gin.Default()

	r.LoadHTMLFiles("templates/index.html", "templates/data.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/submit", func(c *gin.Context) {
		var user User
		if err := c.ShouldBind(&user); err == nil {
			db.Create(&user)
			c.JSON(http.StatusOK, gin.H{"status": "user submitted"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "submission failed"})
		}
	})

	r.GET("/data", func(c *gin.Context) {
		var users []User
		db.Find(&users)

		var usersWithSerial []UserWithSerial
		for i, user := range users {
			usersWithSerial = append(usersWithSerial, UserWithSerial{
				SerialNumber: i + 1,
				User:         user,
			})
		}

		successMessage := c.Query("success")
		c.HTML(http.StatusOK, "data.html", gin.H{
			"users":   usersWithSerial,
			"success": successMessage,
		})
	})

	r.GET("/crash", func(c *gin.Context) {
		panic("Deliberate crash")
	})

	r.Run(":8080")
}
