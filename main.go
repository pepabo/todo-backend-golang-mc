package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>{{.Title}}</title>
    <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
    <link rel="stylesheet" href="https://code.getmdl.io/1.3.0/material.indigo-pink.min.css">
    <script defer src="https://code.getmdl.io/1.3.0/material.min.js"></script>
	</head>
	<body>
    <table class="mdl-data-table mdl-js-data-table mdl-shadow--2dp">
      <thead>
        <tr>
          <th class="mdl-data-table__cell--non-numeric">UA</th>
          <th>Created At</th>
        </tr>
      </thead>
      <tbody>
        {{range $i, $l := .Logs}}
        <tr>
          <td class="mdl-data-table__cell--non-numeric">{{ $l.UA }}</td><td>{{ $l.CreatedAt }}</td>
        </tr>
        {{end}}
      <tbody>
    </table>
	</body>
</html>`

// Access ...
type Access struct {
	UA        string
	CreatedAt time.Time
}

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

		logs := []Access{}
		for rows.Next() {
			var ua string
			var createdAt time.Time

			if err := rows.Scan(&ua, &createdAt); err != nil {
				log.Fatal(err)
			}
			a := Access{
				UA:        ua,
				CreatedAt: createdAt,
			}
			logs = append(logs, a)
		}

		t, err := template.New("page").Parse(tpl)
		if err != nil {
			log.Fatal(err)
		}

		data := struct {
			Title string
			Logs  []Access
		}{
			Title: "Access logs",
			Logs:  logs,
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
