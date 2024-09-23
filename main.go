package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

type Task struct {
	ID          int
	Description string
	CreatedAt   time.Time
	IsComplete  bool
}

var tasksFile = "tasks.csv"

func main() {
	var rootCmd = &cobra.Command{Use: "tasks"}

	var addCmd = &cobra.Command{
		Use:   "add <description>",
		Short: "Add a new task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			addTask(args[0])
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		Run: func(cmd *cobra.Command, args []string) {
			showAll, _ := cmd.Flags().GetBool("all")
			listTasks(showAll)
		},
	}
	listCmd.Flags().BoolP("all", "a", false, "Show all tasks")

	var completeCmd = &cobra.Command{
		Use:   "complete <taskid>",
		Short: "Complete a task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			completeTask(args[0])
		},
	}

	var deleteCmd = &cobra.Command{
		Use:   "delete <taskid>",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			deleteTask(args[0])
		},
	}

	rootCmd.AddCommand(addCmd, listCmd, completeCmd, deleteCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func loadFile() (*os.File, error) {
	f, err := os.OpenFile(tasksFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %v", err)
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, err
	}

	return f, nil
}

func closeFile(f *os.File) error {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return f.Close()
}

func addTask(description string) {
	f, err := loadFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer closeFile(f)

	reader := csv.NewReader(f)
	records, _ := reader.ReadAll()

	var newID int
	if len(records) > 0 {
		lastRecord := records[len(records)-1]
		newID, _ = strconv.Atoi(lastRecord[0])
		newID++
	} else {
		newID = 1
	}

	task := Task{
		ID:          newID,
		Description: description,
		CreatedAt:   time.Now(),
		IsComplete:  false,
	}

	writer := csv.NewWriter(f)
	defer writer.Flush()

	if err := writer.Write([]string{
		strconv.Itoa(task.ID),
		task.Description,
		task.CreatedAt.Format(time.RFC3339),
		strconv.FormatBool(task.IsComplete),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing task: %v\n", err)
	}
	fmt.Printf("Added task: %s\n", task.Description)
}

func listTasks(showAll bool) {
	f, err := loadFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer closeFile(f)

	reader := csv.NewReader(f)
	records, _ := reader.ReadAll()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	if showAll {
		fmt.Fprintln(w, "ID\tTask\tCreated\tDone")
	} else {
		fmt.Fprintln(w, "ID\tTask\tCreated")
	}

	for _, record := range records {
		id, _ := strconv.Atoi(record[0])
		desc := record[1]
		createdAt, _ := time.Parse(time.RFC3339, record[2])
		isComplete, _ := strconv.ParseBool(record[3])

		if !showAll && isComplete {
			continue
		}

		if showAll {
			fmt.Fprintf(w, "%d\t%s\t%s\t%v\n", id, desc, timediff.TimeDiff(createdAt), isComplete)
		} else {
			fmt.Fprintf(w, "%d\t%s\t%s\n", id, desc, timediff.TimeDiff(createdAt))
		}
	}
	w.Flush()
}

func completeTask(taskID string) {
	f, err := loadFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer closeFile(f)

	reader := csv.NewReader(f)
	records, _ := reader.ReadAll()
	for i, record := range records {
		id, _ := strconv.Atoi(record[0])
		if strconv.Itoa(id) == taskID {
			records[i][3] = "true"
			break
		}
	}

	f.Truncate(0)
	f.Seek(0, 0)
	writer := csv.NewWriter(f)
	defer writer.Flush()
	writer.WriteAll(records)
	fmt.Printf("Marked task %s as complete.\n", taskID)
}

func deleteTask(taskID string) {
	f, err := loadFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer closeFile(f)

	reader := csv.NewReader(f)
	records, _ := reader.ReadAll()
	var updatedRecords [][]string

	for _, record := range records {
		id, _ := strconv.Atoi(record[0])
		if strconv.Itoa(id) != taskID {
			updatedRecords = append(updatedRecords, record)
		}
	}

	f.Truncate(0)
	f.Seek(0, 0)
	writer := csv.NewWriter(f)
	defer writer.Flush()
	writer.WriteAll(updatedRecords)
	fmt.Printf("Deleted task %s.\n", taskID)
}
