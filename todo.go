package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// Todo ...
type Todo struct {
	ID        int    `json:"-" db:"id"`
	Title     string `json:"title" db:"title"`
	Completed bool   `json:"completed" db:"completed"`
	Order     int    `json:"order" db:"order"`
	URL       string `json:"url"`
}

// TodoService ...
type TodoService struct {
	db *sqlx.DB
}

// NewTodoService ...
func NewTodoService() *TodoService {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	t := TodoService{
		db: db,
	}
	return &t
}

// GetAll ...
func (t *TodoService) GetAll() ([]*Todo, error) {
	result := []Todo{}
	err := t.db.Select(&result, "SELECT * FROM todos ORDER BY `order` ASC")
	if err != nil {
		return nil, err
	}
	todos := []*Todo{}
	for _, t := range result {
		todos = append(todos, &Todo{
			ID:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
			Order:     t.Order,
		})
	}
	return todos, nil
}

// Get ...
func (t *TodoService) Get(id int) (*Todo, error) {
	todo := Todo{}
	err := t.db.Get(&todo, "SELECT * FROM todos WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// Save ...
func (t *TodoService) Save(todo *Todo) error {
	if todo.ID == 0 {
		tx := t.db.MustBegin()
		result, _ := tx.NamedExec("INSERT INTO todos (title, completed, `order`) VALUES (:title, :completed, :order)", todo)
		lastID, _ := result.LastInsertId()
		todo.ID = int(lastID)
		tx.Commit()
		return nil
	}
	tx := t.db.MustBegin()
	_, err := tx.NamedExec("UPDATE todos SET title = :title, completed = :completed,  `order` = :order WHERE id = :id", todo)
	if err != nil {
		return err
	}
	tx.Commit()
	err = t.db.Get(todo, "SELECT * FROM todos WHERE id = ?", todo.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAll ...
func (t *TodoService) DeleteAll() error {
	tx := t.db.MustBegin()
	tx.MustExec("DELETE FROM todos")
	tx.Commit()
	return nil
}

// Delete ...
func (t *TodoService) Delete(id int) error {
	tx := t.db.MustBegin()
	tx.MustExec("DELETE FROM todos WHERE id = ?", id)
	tx.Commit()
	return nil
}
