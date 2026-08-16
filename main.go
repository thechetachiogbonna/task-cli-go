package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const dataFile = "tasks.json"

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
}

func main() {
	tasks, err := loadTasks(dataFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading tasks %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "add":
		addTask(tasks, args)
	case "list":
		listTasks(tasks)
	case "update":
		updateTask(tasks, args)
	case "delete":
		deleteTask(tasks, args)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q\n", command)
		printUsage()
	}
}

func printUsage() {
	println("Usage:")
	println("  add <task description> - Add a new task")
	println("  list [--completed|--pending|--all] - List all tasks")
	println("  update id=<task ID> [description=<new description>] [done=<true|false>] - Update a task")
	println("  delete id=<task ID> - Delete a task")
	println("  help|-h|--help - Show this help message")
}

func loadTasks(path string) ([]Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}

		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func listTasks(tasks []Task)                 {}
func addTask(tasks []Task, args []string)    {}
func updateTask(tasks []Task, args []string) {}
func deleteTask(tasks []Task, args []string) {}
