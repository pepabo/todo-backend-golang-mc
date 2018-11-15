package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type todoHandler struct {
	service *TodoService
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	handler := todoHandler{
		service: NewTodoService(),
	}

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

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

	log.Println("Start server.")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func (h *todoHandler) listTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addURLToTodos(r, todos...)
	json.NewEncoder(w).Encode(todos)
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
	json.NewEncoder(w).Encode(todo)
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
	json.NewEncoder(w).Encode(todo)
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
	json.NewEncoder(w).Encode(todo)
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
	h.service.DeleteAll()
	w.WriteHeader(http.StatusNoContent)
}

func addURLToTodos(r *http.Request, todos ...*Todo) {
	scheme := "https"
	baseURL := scheme + "://" + r.Host + "/todos/"

	for _, todo := range todos {
		todo.URL = baseURL + strconv.Itoa(todo.ID)
	}
}
