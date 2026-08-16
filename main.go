package main

import "os"

const dataFile = "tasks.json"

func main() {
	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "add":
		addTask(args)
	case "list":
		listTasks()
	case "update":
		updateTask(args)
	case "delete":
		deleteTask(args)
	case "help":
		printUsage()
	default:
		printUsage()
	}
}

func printUsage() {
	println("Usage:")
	println("  add <task description> - Add a new task")
	println("  list [--completed|--pending|--all] - List all tasks")
	println("  update id=<task ID> [description=<new description>] [done=<true|false>] - Update a task")
	println("  delete id=<task ID> - Delete a task")
	println("  help - Show this help message")
}

func listTasks()               {}
func addTask(args []string)    {}
func updateTask(args []string) {}
func deleteTask(args []string) {}
