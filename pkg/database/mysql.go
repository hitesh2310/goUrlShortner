package database

import (
	"database/sql"
	"fmt"
	"log"
	"main/logs"
	"main/pkg/constants"
	"net/url"
	"time"

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
		fmt.Println("Not connected ::", err)
		return
	}

	DbConn = db
	fmt.Println("Connected to the database")

}

func AddEntry(longUrl, shortUrl string) error {

	// query := fmt.Sprintf("INSERT INTO url_mapping (long_url, short_url) VALUES (%s,%s)", longUrl, shortUrl)
	stmt, err := DbConn.Prepare("INSERT INTO url_mapping (longUrl, shortUrl,createdAt) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	currentTime := time.Now()

	epochTime := currentTime.Unix()

	_, err = stmt.Exec(longUrl, shortUrl, epochTime)
	if err != nil {
		fmt.Println("Error executing the query::", err)
		return err
	}

	fmt.Println("User inserted successfully")
	return nil
}

func GetShortUrl(longUrl string) (string, error) {
	var shortUrl string
	query := "SELECT shortUrl FROM url_mapping WHERE longUrl = ?"
	rows, err := DbConn.Query(query, longUrl)
	if err != nil {
		fmt.Println("Error to get shortUrl::", err)
		return "", err
	}
	defer rows.Close()

	// Iterate over the rows and print the shortUrl values
	for rows.Next() {
		if err := rows.Scan(&shortUrl); err != nil {
			fmt.Println("Error to scan short Url::", err)
			return "", err
		}
		fmt.Println("Short URL:", shortUrl)
	}
	return shortUrl, nil
}

func GetCount() (int, error) {

	query := "SELECT COUNT(*) FROM url_mapping"
	var count int
	err := DbConn.QueryRow(query).Scan(&count)
	if err != nil {
		fmt.Println("Error querying ::", err)
	}

	fmt.Println("Row count:", count)
	return count, nil
}

func GetEntry(shortUrl string) string {
	query := "SELECT longUrl FROM url_mapping WHERE shortUrl = ?"
	var longURL string
	err := DbConn.QueryRow(query, shortUrl).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No matching record found")
			return ""
		} else {
			log.Fatal(err)
			return ""
		}
	}

	return longURL
}

func IncrementHost(longURL string) {
	parsedURL, err := url.Parse(longURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	host := parsedURL.Host

	// Prepare the SQL statement
	stmt, err := DbConn.Prepare("INSERT INTO host_stats (host, count) VALUES (?, 1) ON DUPLICATE KEY UPDATE count = count + 1")
	if err != nil {
		fmt.Println("Error to prepare query::", err.Error())
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(host)
	if err != nil {
		fmt.Println("Error to execute the query::", err.Error())
	}
}

func GetTop3() []string {
	fmt.Println("COde check top")
	var result []string
	query := "SELECT host FROM host_stats ORDER BY count DESC LIMIT 3"
	rows, err := DbConn.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var host string
		err := rows.Scan(&host)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Host:", host)
		result = append(result, host)
	}
	return result
}
