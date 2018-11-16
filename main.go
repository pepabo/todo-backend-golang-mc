package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// `/todos` のリクエストを取り扱うハンドラ
type todoHandler struct {
	service *TodoService
}

func main() {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	// 実行ファイルを同じディレクトリの.envファイルを読み込む
	err = godotenv.Load(filepath.Join(filepath.Dir(exe), ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	handler := todoHandler{
		service: NewTodoService(),
	}

	r := chi.NewRouter()

	// CORSの設定 go-chi のミドルウェア機構を使用
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	// `/todos` のルーティング
	r.Route("/todos", func(r chi.Router) {
		r.Get("/", handler.listTodos)
		r.Post("/", handler.createTodo)
		r.Delete("/", handler.deleteAllTodos)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.getTodo)
			r.Patch("/", handler.updateTodo)
			r.Delete("/", handler.deleteTodo)
		})
	})

	// `/` のルーティング
	r.Get("/", indexHandler)

	log.Println("Start server.")
	// Webサーバの起動 ( 0.0.0.0:8080でlisten。ルーティングにgo-chiを使用 )
	log.Fatal(http.ListenAndServe(":8080", r))
}

func (h *todoHandler) listTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addURLToTodos(r, todos...)
	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *todoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	todo := Todo{
		Completed: false,
	}
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	err = h.service.Save(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addURLToTodos(r, &todo)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *todoHandler) getTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}
	todo, err := h.service.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if todo == nil {
		http.NotFound(w, r)
		return
	}
	addURLToTodos(r, todo)
	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *todoHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}
	todo, err := h.service.Get(id)
	if err != nil {
		if strings.ToLower(err.Error()) == "not found" {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewDecoder(r.Body).Decode(todo)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	err = h.service.Save(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addURLToTodos(r, todo)
	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *todoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}
	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *todoHandler) deleteAllTodos(w http.ResponseWriter, r *http.Request) {
	err := h.service.DeleteAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	// テンプレートファイルの読み込み
	t := template.Must(template.ParseFiles(filepath.Join(filepath.Dir(exe), "templates/index.html.tpl")))
	data := struct {
		Title string
	}{
		Title: "Todo-Backend API server written in Go for \"LOLIPOP! Managed Cloud\"",
	}
	// テンプレートのレンダリング
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func addURLToTodos(r *http.Request, todos ...*Todo) {
	scheme := "https"
	baseURL := scheme + "://" + r.Host + "/todos/"

	for _, todo := range todos {
		todo.URL = baseURL + strconv.Itoa(todo.ID)
	}
}
