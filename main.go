package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"time"
)

// Colors for different status
const (
	Green  = "\033[92m"
	Yellow = "\033[93m"
	Blue   = "\033[94m"
	Reset  = "\033[0m"
)

// Initiallizing a slice of tasks that is modeled after the task objects in the JSON file
var tasks []*Task

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created-at"`
	UpdatedAt   string `json:"updated-at"`
}

// Time formatting for readable ouput and storing status as a constant for easy modifications
const (
	customTimeFormat = "2006/1/2 15:04"
	shortTimeFormat  = "2006/01/02 15:04"

	todoStatus       = "todo"
	inProgressStatus = "in-progress"
	doneStatus       = "done"
)

// Lists all tasks provided
func PrintTasks(tasks []*Task) {
	for _, task := range tasks {
		// Parse the createdAt and updatedAt fields (assuming customTimeFormat)
		createdAt, err := time.Parse(customTimeFormat, task.CreatedAt)
		if err != nil {
			fmt.Println("Failed to parse CreatedAt:", err)
			return
		}
		updatedAt, err := time.Parse(customTimeFormat, task.UpdatedAt)
		if err != nil {
			fmt.Println("Failed to parse UpdatedAt:", err)
			return
		}

		// Conditionally apply color to the Status field
		var coloredStatus string
		switch task.Status {
		case todoStatus:
			coloredStatus = Blue + strings.ToUpper(task.Status) + Reset
		case doneStatus:
			coloredStatus = Green + strings.ToUpper(task.Status) + Reset
		case inProgressStatus:
			coloredStatus = Yellow + strings.ToUpper(task.Status) + Reset
		default:
			coloredStatus = task.Status // no color for unknown status
		}

		// Print task information with colored status
		fmt.Printf("ID: %d \n%s %s \nCreatedAt: %v, \nUpdatedAt: %v\n\n",
			task.ID, task.Description, coloredStatus, createdAt.Format(shortTimeFormat), updatedAt.Format(shortTimeFormat))
	}
}

// Lists tasks conditionally according to the status argument
func ListTaks(status string) ([]*Task, error) {
	switch status {
	case todoStatus:
		fmt.Println("TODO")
	case inProgressStatus:
		fmt.Println("INPROGRESS")
	case doneStatus:
		fmt.Println("DONE")
	}
	return nil, nil
}

func main() {
	// taskData is an embedded file
	err := json.Unmarshal(tasksData, &tasks)
	if err != nil {
		fmt.Println("Failed to parse \"tasks.json\":", err)
		return
	}

	PrintTasks(tasks)
}
