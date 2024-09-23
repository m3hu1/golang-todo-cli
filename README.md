# Todo CLI

A simple command-line application for managing tasks using Go. This CLI app allows users to perform CRUD operations on tasks, including adding, listing, completing, and deleting tasks.

## Features

- Add new tasks
- List all tasks (completed and uncompleted)
- Mark tasks as complete
- Delete tasks

## Requirements

- Go 1.16 or higher

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/todo-cli.git
   cd todo-cli
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. (Optional) Build the executable:

   ```bash
   go build -o tasks
   ```

## Usage

You can run the application using the following commands:

### Add a Task

```bash
./tasks add "Your task description"
```

### List Tasks

To list all uncompleted tasks:

```bash
./tasks list
```

To list all tasks (including completed ones):

```bash
./tasks list --all
```

### Complete a Task

To mark a task as complete, use its ID:

```bash
./tasks complete <taskid>
```

### Delete a Task

To delete a task by its ID:

```bash
./tasks delete <taskid>
```

## Example Workflow

1. Add tasks:
   ```bash
   ./tasks add "Buy groceries"
   ./tasks add "Clean the house"
   ```

2. List tasks:
   ```bash
   ./tasks list
   ```

3. Complete a task (e.g., ID 1):
   ```bash
   ./tasks complete 1
   ```

4. List all tasks:
   ```bash
   ./tasks list --all
   ```

5. Delete a task (e.g., ID 2):
   ```bash
   ./tasks delete 2
   ```

## Contributing

Feel free to open issues or submit pull requests. Contributions are welcome!
