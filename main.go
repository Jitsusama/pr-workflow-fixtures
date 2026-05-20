package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: tasks <add|list|done|rm> ...")
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "add":
		return cmdAdd(rest)
	case "list":
		return cmdList()
	case "done":
		return cmdDone(rest)
	case "rm":
		return cmdRemove(rest)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}
