package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"time"
)

func GetConnection() *sql.DB {
	// * DB setting
	godotenv.Load("../.env")
	//sqlUsername := os.Getenv("MYSQLUSERNAME")
	//sqlPassword := os.Getenv("MYSQLPASSWORD")
	//sqlPort := os.Getenv("MYSQLPORT")
	//databaseName := "go-learn"

	sqlConn := "root:dotnet66@tcp(localhost:3306)/go-learn"
	//sqlConn := fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s", sqlUsername, sqlPassword, sqlPort, databaseName)
	//fmt.Println(sqlConn, sqlUsername, sqlPassword, sqlPort, databaseName)

	DB, err := sql.Open("mysql", sqlConn)
	if err != nil {
		panic(err)
	}
	// * close otomatis ketika server berhenti
	//defer DB.Close()

	// * setttingan dasar DB conn pooling
	DB.SetMaxIdleConns(10)                  // jumlah koneksi minimal dibuat
	DB.SetMaxOpenConns(100)                 // jumlah koneksi maksimal dibuat
	DB.SetConnMaxIdleTime(5 * time.Minute)  // berapa lama koneksi akan dihapus ketika tidak digunakan
	DB.SetConnMaxLifetime(60 * time.Minute) // berapa lama koneksi boleh digunakan

	return DB
}
