package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetRootHandler(c *gin.Context) {
	// * response
	c.JSON(http.StatusOK, gin.H{
		"erros":   false,
		"message": "Get testing",
	})
}

func GetTestHandler(c *gin.Context) {

	// * get parameter
	id := c.Param("id")

	// * get query dan misal convert ke int
	page_num, _ := strconv.ParseInt(c.Query("page"), 10, 32)
	size, _ := strconv.ParseInt(c.Query("per_page"), 10, 32)

	c.JSON(http.StatusOK, gin.H{
		"erros":   false,
		"message": "Get testing kedua yang ini",
		"data":    id,
		"page":    page_num,
		"size":    size,
	})
}
