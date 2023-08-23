package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"learn-api-go/database"
	"learn-api-go/models"
	"learn-api-go/utils"
	"net/http"
	"strconv"
	"time"
)

// ! --------------------------- CONST & FUNCTION  ---------------------------

/*
// ! diperhatikan
// * ketika misal ingin diekspor dalam bentuk JSON untuk diperhatikan wajib menggunakan huruf Kapital pada huruf pertama agar data terkait dapat di Ekspos dalam bentuk response JSON API
*/
// * struct for user data response json
type UserData struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Age       int    `json:"age"`
	Status    bool   `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// * struct untuk response data untuk get data user misal untuk data get user
type DataResponseJSON struct {
	TotalData int        `json:"total_data"`
	Page      int64      `json:"page"`
	PerPage   int64      `json:"per_page"`
	User      []UserData `json:"users"`
}

// ! --------------------------- ---- ---------------------------

// ! --------------------------- GET ---------------------------

// ! --------------------------- ---- ---------------------------

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

func GetAllUsers(c *gin.Context) {
	db := database.GetConnection()

	ctx := context.Background()

	// ! MANUAL PAGINATION
	//rows, err := db.QueryContext(ctx, "select * from user")
	// ! ---------------------------------

	// * LANGSUNG QUERY DENGAN PAGINATION
	// * ketika query langsung dengan
	page, _ := strconv.ParseInt(c.Query("page"), 10, 32)
	per_page, _ := strconv.ParseInt(c.Query("size"), 10, 32)
	page, per_page = utils.PaginationPageSizeCheck(page, per_page)
	start_data := (page - 1) * per_page
	rows, err := db.QueryContext(ctx, "SELECT * FROM user LIMIT ? OFFSET ?", per_page, start_data) // * harus urut -> limit dahulu baru offset, kebalik akan error
	var total_data int = 10
	err = db.QueryRowContext(ctx, "select count(*) from user").Scan(&total_data)
	// * ---------------------------------
	if err != nil {
		utils.ThrowErr(c, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	// * iterasikan data
	var all_users []UserData
	for rows.Next() {
		var id, age int
		var username, createdAt, updatedAt string
		var status bool
		rows.Scan(&id, &username, &age, &status, &createdAt, &updatedAt)

		user := UserData{
			id, username, age, status, createdAt, updatedAt,
		}
		fmt.Println(user)

		all_users = append(all_users, user)
	}
	// ! TAMBAHAN QUERY PAGINATION
	// * ketika gunakan PAGINATION langsung di query, ketika tidak ada data ditampilkan akan kembalikan (nil) -> perlu konversi untuk kembalikan array kosong
	if all_users == nil {
		all_users = []UserData{}
	}
	// ! -------------------------

	// fmt.Println(all_users)
	// * slicing slice untuk pagination
	// * GET semua data sql dan SLICING MANUAL
	/*
		! catatam MANUAL SLICING
		- ketika slicing dilakukan secara manual maka harus ada beberapa konfigurasi terakit proses slicing slice agar tidak error -> karena ketika misal langsung dislice seperti biasa akan ada kemungkinan untuk parameter misal akhir index data > dari jumlah total data
		- hal di atas tidak bisa diakomodasi (langsung kembalikan slice kosong misalnya) oleh GOlang dan akan ada error
		- untuk mengatasi hal di atas maka proses di bawah akan ada beberapa konfigurasi terlebih dahulu untuk pengecekan dan penyesuaian index untuk slicing
		- KEKURANGAN -> tentu jelas kekurangan ini adalah selain lebih panjang/ ada beberapa tambahan hal di atas juga ketika misal data banyak maka query dulu manual semua data juga akan sangat banyak
		- sedikit kelebihan hanya mungkin jelas query sangat singkat dan satu query tersebut juga secara sekaligus bisa sekalian untuk hitung total data yang ada
		! --------------

		* catatan LANGSUNG QUERY PAGINATION
		- gaada setting misal seperti manual dan data langsung iterasi saja
		- sedikit kekurangan mungkin misal ingin hitung total data berarti harus ada 2 query yaitu 1 query untuk get data sesuai pagination 1 query lagi untuk hitung total data
		* ---------------------------------
	*/
	// * get query parameter for pagination
	//page, _ := strconv.ParseInt(c.Query("page"), 10, 32)
	//per_page, _ := strconv.ParseInt(c.Query("size"), 10, 32)
	//page, per_page = utils.PaginationPageSizeCheck(page, per_page)
	//total_data := len(all_users)
	//start_data := (page - 1) * per_page
	//end_data := start_data + per_page
	//
	//if int(start_data) >= len(all_users) {
	//	all_users = []UserData{}
	//} else {
	//	if int(end_data) > len(all_users) {
	//		end_data = int64(len(all_users))
	//	}
	//
	//	all_users = all_users[start_data:end_data]
	//}

	//all_users = DataResponseJSON{TotalData: len(all_users), Page: page, PerPage: per_page, []UserData{all_users...}}

	c.JSON(http.StatusOK, gin.H{
		"errors":     false,
		"message":    "Get all users data",
		"data":       all_users,
		"per_page":   per_page,
		"page":       page,
		"total_data": total_data,
	})

}

func GetOneUserData(c *gin.Context) {
	id_user := c.Param("id_user")

	db := database.GetConnection()

	ctx := context.Background()

	// * get 1 data user dengan id terkait
	var id, age int
	var username, createdAt, updatedAt string
	var status bool
	idUser, err := strconv.ParseInt(id_user, 10, 32)

	// * PROSES GET 1 DATA/ BARIS DENGAN SQL dan HANDLING ERRONYA
	err = db.QueryRowContext(ctx, "select * from user where id= ?", idUser).Scan(&id, &username, &age, &status, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ThrowErrWithMessage(c, http.StatusBadRequest, "Data tidak ditemukan")
		}
		return
	}
	user := UserData{
		id, username, age, status, createdAt, updatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"errors":  false,
		"message": "Get one data user",
		"data":    user,
	})

}

// ! --------------------------- ---- ---------------------------

// ! --------------------------- POST ---------------------------

// ! --------------------------- ---- ---------------------------

func PostUser(c *gin.Context) {
	var userInput models.User

	err := c.ShouldBindJSON(&userInput)
	if err != nil {
		utils.ThrowErrorValidationJSON(c, http.StatusBadRequest, err)
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

func CreateUserTest(c *gin.Context) {
	db := database.GetConnection()
	defer db.Close() // tutup koneksi kembali di akhir

	// * get user input from post req
	var bodyReq models.User
	// * validation
	errr := c.ShouldBindJSON(&bodyReq)
	if errr != nil {
		utils.ThrowErr(c, http.StatusBadRequest, errr)
		return
	}

	username := bodyReq.Username
	age := bodyReq.Age
	status := bodyReq.Status
	timeNow := time.Now()

	ctx := context.Background()
	_, err := db.ExecContext(ctx, "INSERT INTO user(username, age, status, createdAt, updatedAt) VALUES (?,?,?,?,?)", username, age, status, timeNow, timeNow)
	//_, err := db.ExecContext(ctx, "INSERT INTO user(username, age, status) VALUES ('user1', 12, false)")

	if err != nil {
		utils.ThrowErrWithMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errors":  false,
		"message": "Berhasil input data baru",
	})
}

// ! --------------------------- ---- ---------------------------

// ! --------------------------- UPDATE ---------------------------

// ! --------------------------- ---- ---------------------------

func UpdateUserData(c *gin.Context) {
	id_user := c.Param("id_user")
	idUser, _ := strconv.ParseInt(id_user, 10, 32)

	db := database.GetConnection()
	ctx := context.Background()

	// ! --------------- MISAL CEK DATA INGIN DI EDIT ADA TIDAK DATANYA
	// * karena hanya checking maka hanya cek id saja sudah cukup
	var id int

	// * MISAL
	err := db.QueryRowContext(ctx, "select id from user where id= ?", idUser).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ThrowErrWithMessage(c, http.StatusNotFound, "Data tidak ditemukan, proses edit gagal")
		}
		return
	}
	// ! --------------- --------------- --------------- ---------------

	// * --------------- --------------- user ditemukan -> proses edit data terkait
	// * cek misal ada data tidak sesuai dari input json yang diberikan
	var bodyReq models.User
	err = c.ShouldBindJSON(&bodyReq)
	if err != nil {
		utils.ThrowErrorValidationJSON(c, http.StatusBadRequest, err)
		return
	}
	// * --------------- --------------- --------------- ---------------
	// * Proses update di query
	username := bodyReq.Username
	age := bodyReq.Age
	status := bodyReq.Status
	_, err = db.ExecContext(ctx, "update user set username = ?, age = ?, status = ? where id = ?", username, age, status, idUser)

	// ketika query gagal
	if err != nil {
		utils.ThrowErrWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errors":  false,
		"message": "Berhasil update data",
	})
}

// ! --------------------------- ---- ---------------------------

// ! --------------------------- DELETE ---------------------------

// ! --------------------------- ---- ---------------------------

func DeleteUserData(c *gin.Context) {
	id_user := c.Param("id_user")
	idUser, _ := strconv.ParseInt(id_user, 10, 31)

	db := database.GetConnection()
	ctx := context.Background()

	// ! --------------- MISAL CEK DATA INGIN DI EDIT ADA TIDAK DATANYA
	var id int
	err := db.QueryRowContext(ctx, "select id from user where id = ?", idUser).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ThrowErrWithMessage(c, http.StatusNotFound, "Data tidak ditemukan, proses hapus gagal")
		}
		return
	}
	// ! --------------- --------------- ------------------------------

	// * data ada maka proses hapus
	_, err = db.ExecContext(ctx, "delete from user where id = ?", idUser)
	if err != nil {
		utils.ThrowErrWithMessage(c, http.StatusInternalServerError, "Terjadi error saat proses hapus data")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errors":  false,
		"message": "Berhasil delete data",
	})
}
