package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"time"
)

const help string = `Available commands: 
'list': 
	tasky list (Lists all tasks)
	tasky list <todo|done|in-progress> (Lists tasks with the specified status)
'add': 
	tasky add <task description> (adds a new task)
'update': 
	tasky update <task id> <new description> (update the task description of the task with the specified id)
'delete': 
	tasky delete <task id> (deletes the task with the specified id)
'clear':
	tasky clear (deletes all tasks)
'doing':
	tasky doing <task id> (Assigns the status 'in-progress' to the task with the specified id)
'done':
	tasky done <task id> (Assigns the status 'done' to the task with the specified id)`

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
func listTasks(tasks []*Task) {
	if len(tasks) == 0 {
		fmt.Println("You don't have any tasks")
		return
	}
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

func filterTasks(tasks []*Task, status string) []*Task {
	var output []*Task
	for _, task := range tasks {
		if task.Status == status {
			output = append(output, task)
		}
	}
	return output
}

func addTask(description string) {
	task := Task{
		ID:          len(tasks) + 1,
		Description: description,
		Status:      todoStatus,
		CreatedAt:   time.Now().Format(customTimeFormat),
		UpdatedAt:   time.Now().Format(customTimeFormat),
	}
	tasks = append(tasks, &task)

	err := saveTask("./tasks.json", tasks)
	if err != nil {
		fmt.Println("Failed to save task:", err)
		return
	}

	fmt.Println("Task added successfully!")
}

func updateTask(idInput string, description string) {
	id, err := strconv.Atoi(idInput)
	if err != nil {
		fmt.Printf("Error reading the ID: %s \nhint: ID has to be an integer", idInput)
		return
	}
	for _, task := range tasks {
		if id == task.ID {
			task.Description = description
			task.UpdatedAt = time.Now().Format(customTimeFormat)
			break
		}
	}

	err = saveTask("./tasks.json", tasks)
	if err != nil {
		fmt.Println("Failed to update task:", err)
		return
	}

	fmt.Println("Task updated successfully!")
}

func deleteTask(idInput string) {
	id, err := strconv.Atoi(idInput)
	if err != nil {
		fmt.Printf("Error reading the ID: %s \nID has to be an integer", idInput)
		return
	}
	for i, task := range tasks {
		if id == task.ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	err = saveTask("./tasks.json", tasks)
	if err != nil {
		fmt.Println("Failed to delete task:", err)
		return
	}

	fmt.Println("Task deleted successfully!")
}

func deleteAllTasks() {
	tasks = tasks[:0]

	err := saveTask("./tasks.json", tasks)
	if err != nil {
		fmt.Println("Failed to delete tasks:", err)
		return
	}

	fmt.Println("Tasks deleted successfully!")
}

func assignStatus(idInput string, status string) {
	id, err := strconv.Atoi(idInput)
	if err != nil {
		fmt.Printf("Error reading the ID: %s \nID has to be an integer", idInput)
		return
	}
	for _, task := range tasks {
		if id == task.ID {
			task.Status = status
			task.UpdatedAt = time.Now().Format(customTimeFormat)
			break
		}
	}

	err = saveTask("./tasks.json", tasks)
	if err != nil {
		fmt.Println("Failed to update task:", err)
		return
	}

	fmt.Println("Task status successfully udpated to 'done'!")
}

func saveTask(filePath string, tasks []*Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func handleFile() ([]byte, error) {
	file, err := os.OpenFile("./tasks.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size() == 0 {
		_, err = file.Write([]byte("[]"))
		if err != nil {
			fmt.Println("Error initializing file with an empty array:", err)
			return nil, err
		}
		fmt.Println("Initialized the file with an empty array.")
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func main() {
	bytes, err := handleFile()
	if err != nil {
		fmt.Println("Failed to read or create \"tasks.json\":", err)
		return
	}

	err = json.Unmarshal(bytes, &tasks)
	if err != nil {
		fmt.Println("Failed to parse \"tasks.json\":", err)
		fmt.Println("Try running the app again")
		return
	}

	if (len(os.Args)) < 2 {
		fmt.Println(help)
	}

	command := os.Args[1]

	switch command {
	case "list":
		if len(os.Args) < 3 {
			listTasks(tasks)
		} else {
			status := os.Args[2]
			switch status {
			case "todo":
				listTasks(filterTasks(tasks, todoStatus))
			case "done":
				listTasks(filterTasks(tasks, doneStatus))
			case "doing":
				listTasks(filterTasks(tasks, inProgressStatus))
			default:
				fmt.Printf("Unknown status: %s \nUsage of 'list':\ntasky list (Lists all tasks)\ntasky list <todo|done|in-progress> (lists tasks with the specified status)", status)
			}
		}
	case "add":
		if len(os.Args) < 3 {
			fmt.Printf("Usage of 'add': tasky add <task description>")
		} else {
			description := strings.Join(os.Args[2:], " ")
			addTask(description)
		}
	case "update":
		if len(os.Args) < 4 {
			fmt.Printf("Usage of 'update': tasky update <task id> <new task description>")
		} else {
			description := strings.Join(os.Args[3:], " ")
			updateTask(os.Args[2], description)
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Printf("Usage of 'delete': tasky delete <task ID>")
		} else {
			deleteTask(os.Args[2])
		}
	case "clear":
		deleteAllTasks()
	case "done":
		if len(os.Args) < 3 {
			fmt.Printf("Usage of 'done': tasky done <task ID>")
		} else {
			assignStatus(os.Args[2], doneStatus)
		}
	case "doing":
		if len(os.Args) < 3 {
			fmt.Printf("Usage of 'done': tasky done <task ID>")
		} else {
			assignStatus(os.Args[2], inProgressStatus)
		}
	default:
		fmt.Println("Unknown command: ", command)
		fmt.Println(help)
	}
}
