package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "add":
		tasks, err = addTask(tasks, args)
	case "list":
		err = listTasks(tasks, args)
	case "update":
		tasks, err = updateTask(tasks, args)
	case "delete":
		tasks, err = deleteTask(tasks, args)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q\n", command)
		printUsage()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		printUsage()
		os.Exit(1)
	}

	saveTask(tasks, dataFile)
}

func printUsage() {
	println("Usage:")
	println("  add <task description> - Add a new task")
	println("  list [--completed|--pending|--all] - List all tasks")
	println("  update id=<task ID> [description=<new description>] [done=<true|false>] - Update a task")
	println("  delete id=<task ID> - Delete a task")
	println("  help|-h|--help - Show this help message")
}

func saveTask(tasks []Task, path string) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
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

func listTasks(tasks []Task, args []string) error {
	var filter string

	if len(args) > 0 {
		filter = args[0]
	}

	if filter == "-a" {
		filter = "--all"
	}

	if filter != "--completed" && filter != "--pending" && filter != "--all" {
		return fmt.Errorf("Unknown filter %q\n", filter)
	}

	var filteredTasks []Task
	switch filter {
	case "--completed":
		for _, task := range tasks {
			if task.Done {
				filteredTasks = append(filteredTasks, task)
			}
		}
	case "--pending":
		for _, task := range tasks {
			if !task.Done {
				filteredTasks = append(filteredTasks, task)
			}
		}
	default:
		filteredTasks = tasks
	}

	for _, task := range filteredTasks {
		status := "×"
		if task.Done {
			status = "🗸"
		}
		fmt.Printf("[%s] %d: %s\n", status, task.ID, task.Description)
	}
	return nil
}

func addTask(tasks []Task, args []string) ([]Task, error) {
	if len(args) == 0 {
		return nil, errors.New("Error: Task description is required")
	}

	description := args[0]
	task := Task{
		ID:          len(tasks) + 1,
		Description: description,
		Done:        false,
		CreatedAt:   time.Now(),
	}

	return append(tasks, task), nil
}

func updateTask(tasks []Task, args []string) ([]Task, error) {
	parsedArgs, err := parseArgs(args)
	if err != nil {
		return tasks, err
	}

	id, ok := parsedArgs["id"]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: Task ID is required\n")
		return tasks, errors.New("Task ID is required")
	}

	taskID, err := strconv.Atoi(id)
	if err != nil {
		return tasks, errors.New("Invalid task ID")
	}

	var task *Task
	for _, t := range tasks {
		if t.ID == taskID {
			task = &t
			break
		}
	}

	if task == nil {
		return tasks, errors.New("Task not found")
	}

	description, userProvidedDescription := parsedArgs["description"]
	if userProvidedDescription {
		task.Description = description
	}

	done, userProvidedDone := parsedArgs["done"]
	if userProvidedDone {
		task.Done = done == "true"
	}
	if userProvidedDone && done != "true" && done != "false" {
		return tasks, errors.New("Invalid value for done. Use true or false")
	}

	if !userProvidedDescription && !userProvidedDone {
		return tasks, errors.New("No fields to update")
	}

	tasks[taskID-1] = *task

	return tasks, nil
}

func deleteTask(tasks []Task, args []string) ([]Task, error) {
	parsedArgs, err := parseArgs(args)

	if err != nil {
		return nil, err
	}

	id, ok := parsedArgs["id"]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: Task ID is required\n")
		return tasks, errors.New("Task ID is required")
	}

	taskID, err := strconv.Atoi(id)
	if err != nil {
		return tasks, errors.New("Invalid task ID")
	}

	var task *Task
	var taskIndex int
	for i, t := range tasks {
		if t.ID == taskID {
			task = &t
			taskIndex = i
			break
		}
	}

	if taskIndex < 0 {
		return nil, fmt.Errorf("task with id %d not found", taskID)
	}

	if task == nil {
		return tasks, errors.New("Task not found")
	}

	tasks = append(tasks[:taskIndex], tasks[taskIndex+1:]...)

	return tasks, nil
}

func parseArgs(args []string) (map[string]string, error) {
	parsedArgs := make(map[string]string)
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid argument format: %q", arg)
		}
		parsedArgs[parts[0]] = parts[1]
	}
	return parsedArgs, nil
}
