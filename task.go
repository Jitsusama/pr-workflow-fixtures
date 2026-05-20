package main

import "time"

// Task is a single todo entry persisted to the on-disk
// store.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	Priority  string    `json:"priority,omitempty"`
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

// priorityRank gives a sort weight per priority. Lower is
// more important.
var priorityRank = map[string]int{
	"high":   0,
	"normal": 1,
	"low":    2,
}
