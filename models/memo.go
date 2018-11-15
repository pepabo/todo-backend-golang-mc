package models

import (
	"database/sql"
	"log"
	"time"
)

// Memo ...
type Memo struct {
	Title     string
	Body      string
	CreatedAt time.Time
	UpdateAt  time.Time
}

// GetMemos ...
func GetMemos(conn *sql.DB) ([]Memo, error) {
	rows, err := conn.Query("SELECT title, body, created_at, updated_at FROM memo ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memos := []Memo{}
	for rows.Next() {
		var (
			title     string
			body      string
			createdAt time.Time
			updatedAt time.Time
		)
		if err := rows.Scan(&title, &body, &createdAt, &updatedAt); err != nil {
			log.Fatal(err)
		}
		m := Memo{
			Title:     title,
			Body:      body,
			CreatedAt: createdAt,
			UpdateAt:  updatedAt,
		}
		memos = append(memos, m)
	}

	return memos, nil
}
