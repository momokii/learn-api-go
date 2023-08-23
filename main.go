package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"learn-api-go/controllers"
)

func ThrowError(c *gin.Context, message string, code int) {
	c.JSON(code, gin.H{
		"errors":  true,
		"message": message,
	})
}

func main() {
	cors := cors.Default() // default set
	// * alternatif cors setting
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

	// * use middleware
	router.Use(cors)

	// * grouping router misal jadi /user -> untuk banyak routing
	testGroupRouting := router.Group("/testing")
	postGroupingRouting := router.Group("/postGrouping")

	// * use group routing
	router.GET("/", controllers.GetRootHandler)
	testGroupRouting.GET("/test/:id", controllers.GetTestHandler)
	testGroupRouting.GET("/users/:id_user", controllers.GetOneUserData)
	testGroupRouting.GET("/users", controllers.GetAllUsers)
	postGroupingRouting.POST("/users", controllers.PostUser)
	postGroupingRouting.POST("/userReal", controllers.CreateUserTest)
	router.PATCH("/users/:id_user", controllers.UpdateUserData)
	router.DELETE("/users/:id_user", controllers.DeleteUserData)

	router.Run("localhost:8000")
}
