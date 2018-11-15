package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/k1LoW/mc-go-server/models"
)

var conn *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Get("/", getMemos)
	r.Post("/", postMemo)

	log.Println("Start server.")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getMemos(w http.ResponseWriter, r *http.Request) {
	memos, err := models.GetMemos(conn)
	if err != nil {
		log.Fatal(err)
	}

	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	t := template.Must(template.ParseFiles(filepath.Join(filepath.Dir(exe), "templates/index.html.tpl")))
	data := struct {
		Title string
		Memos []models.Memo
	}{
		Title: "Memo",
		Memos: memos,
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func postMemo(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	values, _ := url.ParseQuery(string(body))

	err := models.CreateMemo(conn)
}
