package main

import (
	"strings"
	"testing"
)

func TestNextID(t *testing.T) {
	cases := []struct {
		name  string
		tasks []Task
		want  int
	}{
		{"empty", []Task{}, 1},
		{"one task", []Task{{ID: 1}}, 2},
		{"with gap", []Task{{ID: 1}, {ID: 3}}, 4},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := nextID(c.tasks); got != c.want {
				t.Errorf("nextID = %d, want %d", got, c.want)
			}
		})
	}
}

func TestNewTaskFields(t *testing.T) {
	tk := NewTask(7, "buy milk")
	if tk.ID != 7 {
		t.Errorf("ID = %d, want 7", tk.ID)
	}
	if tk.Title != "buy milk" {
		t.Errorf("Title = %q, want %q", tk.Title, "buy milk")
	}
	if tk.Done {
		t.Error("Done should default to false")
	}
	if tk.CreatedAt.IsZero() {
		t.Error("CreatedAt should be stamped")
	}
}

// helper: assert error message contains substring
func errContains(err error, sub string) bool {
	return err != nil && strings.Contains(err.Error(), sub)
}

var _ = errContains // keep linkable for command tests
