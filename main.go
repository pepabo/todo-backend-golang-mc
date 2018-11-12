package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")))
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Exec("INSERT INTO access_logs SET ua = ?", r.Header.Get("User-Agent"))
		if err != nil {
			log.Fatal(err)
		}

		rows, err := conn.Query("SELECT ua, created_at FROM access_logs ORDER BY created_at DESC")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var ua string
			var createdAt time.Time

			if err := rows.Scan(&ua, &createdAt); err != nil {
				log.Fatal(err)
			}
			fmt.Fprintln(w, fmt.Sprintf("%s, %v\n", ua, createdAt))
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
