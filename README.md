# Task CLI

A small command-line task manager written in Go. Tasks are stored locally in `tasks.json`, so the application does not require a database or external service.

## Requirements

- Go 1.26.3 or newer

## Build

```bash
go build -o task-cli
```

On Windows, this produces `task-cli.exe`.

## Usage

Run the application with:

```bash
go run . <command> [arguments]
```

### Add a task

```bash
go run . add "Buy groceries"
```

The task starts as pending and receives the next numeric ID.

### List tasks

```bash
go run . list --all
go run . list --pending
go run . list --completed
```

`-a` is also accepted as an alias for `--all`.

### Update a task

Arguments use the `key=value` format. You can update the description, completion state, or both:

```bash
go run . update id=1 description="Buy groceries and coffee"
go run . update id=1 done=true
go run . update id=1 description="Buy groceries" done=false
```

The `done` value must be either `true` or `false`.

### Delete a task

```bash
go run . delete id=1
```

### Show help

```bash
go run . help
```

`-h` and `--help` are also supported.

## Data format

Tasks are saved as an indented JSON array in `tasks.json`:

```json
[
  {
    "id": 1,
    "description": "Buy groceries",
    "done": false,
    "created_at": "2026-08-20T12:00:00Z"
  }
]
```

If `tasks.json` does not exist, the CLI starts with an empty task list and creates the file after a successful mutating command.

## Project structure

```text
.
├── main.go     # CLI implementation
├── tasks.json  # Local task storage
└── go.mod      # Go module metadata
```
