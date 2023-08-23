package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// * --------------------------- FUNCTION ---------------------------

// ! ------------- GLOBAL FUNCTION ERROR RETURN

func ThrowErr(c *gin.Context, statusCode int, err error) {
	errorMessages := []string{}
	for _, e := range err.(validator.ValidationErrors) {
		errorMessage := fmt.Sprintf("Error pada field %s dimana seharusnya %s", e.Field(), e.ActualTag())
		errorMessages = append(errorMessages, errorMessage)
	}

	c.JSON(statusCode, gin.H{
		"errors":  true,
		"message": errorMessages,
	})
}

func ThrowErrWithMessage(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"errors":  true,
		"message": message,
	})
}

func ThrowErrorValidationJSON(c *gin.Context, statusCode int, err error) {
	errorMessages := []string{}
	for _, e := range err.(validator.ValidationErrors) {
		errorMessage := fmt.Sprintf("Error ada di field %s, kesalahan pada field tersebut: %s", e.Field(), e.ActualTag())
		errorMessages = append(errorMessages, errorMessage)
	}
	c.JSON(statusCode, gin.H{
		"errors":  true,
		"message": errorMessages,
	})

}

// ! ------------- PAGINATION CHECKING PAGE & SIZE

// * gunakan int6 karena menyesuaikan konversi string ke int -> strconv.ParseInt(c.Query("size"), 10, 32)
func PaginationPageSizeCheck(page, per_page int64) (int64, int64) {
	if page == 0 || per_page == 0 {
		if page == 0 {
			page = 1
		}
		if per_page == 0 {
			per_page = 10
		}
	}
	return page, per_page
}
