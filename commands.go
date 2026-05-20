package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// cmdAdd appends a new task to the store. The title is the
// concatenation of all positional args; an optional
// --priority <level> sets the task's priority.
func cmdAdd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: tasks add <title> [--priority <level>]")
	}

	priority := ""
	var titleParts []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--priority" && i+1 < len(args) {
			priority = args[i+1]
			i++
			continue
		}
		titleParts = append(titleParts, args[i])
	}
	title := strings.Join(titleParts, " ")

	tasks, err := Load()
	if err != nil {
		return err
	}
	id := nextID(tasks)
	task := NewTask(id, title)
	task.Priority = priority
	tasks = append(tasks, task)
	if err := Save(tasks); err != nil {
		return err
	}
	fmt.Printf("added task %d: %s\n", id, title)
	return nil
}

// cmdList prints every task with a done-marker. Pass
// --sort priority to order high-to-low.
func cmdList(args []string) error {
	tasks, err := Load()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("no tasks")
		return nil
	}

	for i := 0; i < len(args); i++ {
		if args[i] == "--sort" && i+1 < len(args) && args[i+1] == "priority" {
			sort.Slice(tasks, func(a, b int) bool {
				return priorityRank[tasks[a].Priority] < priorityRank[tasks[b].Priority]
			})
			i++
		}
	}

	for _, t := range tasks {
		marker := " "
		if t.Done {
			marker = "x"
		}
		if t.Priority != "" {
			fmt.Printf("[%s] %d. (%s) %s\n", marker, t.ID, t.Priority, t.Title)
		} else {
			fmt.Printf("[%s] %d. %s\n", marker, t.ID, t.Title)
		}
	}
	return nil
}

// cmdDone marks the task with the given id as done.
func cmdDone(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: tasks done <id>")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid id: %s", args[0])
	}
	tasks, err := Load()
	if err != nil {
		return err
	}
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Done = true
			return Save(tasks)
		}
	}
	return fmt.Errorf("no task with id %d", id)
}

// cmdRemove drops the task with the given id from the store.
func cmdRemove(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: tasks rm <id>")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid id: %s", args[0])
	}
	tasks, err := Load()
	if err != nil {
		return err
	}
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return Save(tasks)
		}
	}
	return fmt.Errorf("no task with id %d", id)
}

func nextID(tasks []Task) int {
	id := 1
	for _, t := range tasks {
		if t.ID >= id {
			id = t.ID + 1
		}
	}
	return id
}
