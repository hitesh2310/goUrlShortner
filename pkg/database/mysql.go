package database

import (
	"database/sql"
	"fmt"
	"main/logs"
	"main/pkg/constants"

	_ "github.com/go-sql-driver/mysql"
)

var DbConn *sql.DB

func EstablishDbConnection() {
	fmt.Println("Establishing DB connection")
	host := constants.ApplicationConfig.Database.Host
	port := constants.ApplicationConfig.Database.Port
	username := constants.ApplicationConfig.Database.Username
	password := constants.ApplicationConfig.Database.Password
	dbname := constants.ApplicationConfig.Database.Dbname
	// fmt.Println("HEREE")
	// fmt.Println(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	fmt.Println(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	if err != nil {
		logs.InfoLog(err.Error())
	}
	// defer db.Close()
	fmt.Println("Conn opened")
	// Check the database connection
	fmt.Println("Ping check ")
	err = db.Ping()
	fmt.Println("Ping checked")

	if err != nil {
		// logs.InfoLog(err.Error())
		fmt.Println("Not connected %v", err)
		return
	}

	DbConn = db
	fmt.Println("Connected to the database")

}
