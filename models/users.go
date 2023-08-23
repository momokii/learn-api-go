package models

import "encoding/json"

type User struct {
	Username string      `json:"username" binding:"required"` // * gunakan backtics bukan '' dan tidak boleh ada spasi antara misal json:"value1,value2"
	Age      json.Number `json:"age" binding:"required,number"`
	Status   bool        `json:"status"` // ini tidak bisa di set required ->
}
