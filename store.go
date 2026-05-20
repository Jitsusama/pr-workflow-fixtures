package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const storeFileName = ".tasks.json"

func storePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, storeFileName), nil
}

// Load reads the on-disk task list. Returns an empty slice
// if the store doesn't exist yet.
func Load() ([]Task, error) {
	path, err := storePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Save writes the task list back to disk, atomically
// replacing the previous contents.
func Save(tasks []Task) error {
	path, err := storePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
