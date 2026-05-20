package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// validPriorities lists the priority levels cmdAdd accepts.
// Kept in sync with priorityRank in task.go.
var validPriorities = []string{"high", "normal", "low"}

func isValidPriority(p string) bool {
	for _, v := range validPriorities {
		if p == v {
			return true
		}
	}
	return false
}

// cmdAdd appends a new task to the store. The title is the
// concatenation of all positional args; an optional
// --priority <level> sets the task's priority.
func cmdAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	priority := fs.String("priority", "", fmt.Sprintf("task priority (%s)", strings.Join(validPriorities, "|")))
	usage := "usage: tasks add <title> [--priority <level>]"
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("%s: %w", usage, err)
	}

	title := strings.Join(fs.Args(), " ")
	if title == "" {
		return errors.New(usage)
	}
	if *priority != "" && !isValidPriority(*priority) {
		return fmt.Errorf("invalid priority %q: must be one of %s", *priority, strings.Join(validPriorities, ", "))
	}

	// Tag parsing: scan args by hand for --tag <value>.
	tag := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--tag" && i+1 < len(args) {
			tag = args[i+1]
			i++
		}
	}

	tasks, err := Load()
	if err != nil {
		return err
	}
	id := nextID(tasks)
	task := NewTask(id, title)
	task.Priority = *priority
	task.Tag = tag
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
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	sortBy := fs.String("sort", "", "sort order (priority)")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("usage: tasks list [--sort priority]: %w", err)
	}

	tasks, err := Load()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("no tasks")
		return nil
	}

	switch *sortBy {
	case "":
		// keep insertion order
	case "priority":
		sort.Slice(tasks, func(a, b int) bool {
			return priorityRank[tasks[a].Priority] < priorityRank[tasks[b].Priority]
		})
	default:
		return fmt.Errorf("unknown sort order %q: must be priority", *sortBy)
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
