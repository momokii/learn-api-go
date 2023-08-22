package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"learn-api-go/controllers"
	"net/http"
)

type Response struct {
	Errors  bool
	Message string
}

func ThrowError(c *gin.Context, message string, code int) {
	c.JSON(code, gin.H{
		"errors":  true,
		"message": message,
	})
}

func main() {
	cors := cors.Default() // default set
	//cors := cors.New(cors.Config{
	//	AllowOrigins:     []string{"https://foo.com"},
	//	AllowMethods:     []string{"PUT", "PATCH"},
	//	AllowHeaders:     []string{"Origin"},
	//	ExposeHeaders:    []string{"Content-Length"},
	//	AllowCredentials: true,
	//	AllowOriginFunc: func(origin string) bool {
	//		return origin == "https://github.com"
	//	},
	//	MaxAge: 12 * time.Hour,
	// })

	router := gin.Default()

	// * middleware
	router.Use(cors)

	// * grouping router misal jadi /user -> untuk banyak routing

	testGroupRouting := router.Group("/testing")
	postGroupingRouting := router.Group("/postGrouping")

	router.GET("/", controllers.GetRootHandler)
	testGroupRouting.GET("/test/:id", controllers.GetTestHandler)
	postGroupingRouting.POST("/user", PostUser)

	router.Run("localhost:8000")
}

type User struct {
	Username string      `json:"username" binding:"required"` // * gunakan backtics bukan '' dan tidak boleh ada spasi antara misa json:"value1,value2"
	Age      json.Number `json:"age" binding:"required,number"`
	Status   bool        `json:"status"`
}

func PostUser(c *gin.Context) {
	var userInput User

	err := c.ShouldBindJSON(&userInput)
	if err != nil {
		// * loop error validation
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error ada di field %s, karena %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":   true,
			"messages": errorMessages,
		})
		return

		// * jadi throw error di buat func biasa seperti di node
		//ThrowError(c, "Ada Error", 400)
		//return
		//c.JSON(http.StatusBadRequest, err)
		//fmt.Println(err)
		//return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "success add user",
		"username": userInput.Username,
		"umur":     userInput.Age,
		"status":   userInput.Status,
	})
}
