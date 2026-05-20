package main

import "time"

// Task is a single todo entry persisted to the on-disk
// store.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// NewTask stamps a new task with the given id and title.
func NewTask(id int, title string) Task {
	return Task{
		ID:        id,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
}
