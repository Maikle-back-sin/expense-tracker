package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	FileName         = "task.json"
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
	StatusDone       = "done"
)

type Purchase struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Date        time.Time `json:"created_at"`
	Amount      int       `json:"amount"`
}

func loadTasks() ([]Purchase, error) {
	var purchases []Purchase
	file, err := os.Open(FileName)
	if err != nil {
		if os.IsNotExist(err) {
			return purchases, nil
		}
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&purchases)
	return purchases, err
}

func saveTasks(purchases []Purchase) error {
	file, err := os.Create(FileName)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(purchases)
}

func nextID(purchases []Purchase) int {
	max := 0
	for _, t := range purchases {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

func findTask(purchases []Purchase, id int) (*Purchase, int) {
	for i, t := range purchases {
		if t.ID == id {
			return &purchases[i], i
		}
	}
	return nil, -1
}

func cmdAdd(args []string) error {
	if len(args) < 5 {
		return fmt.Errorf("usage: add \"--description *Your desc* --amount *Price*\"")
	}
	description := args[2]
	amount, err := strconv.Atoi(args[4])

	fmt.Println(description, amount)

	purchases, err := loadTasks()
	if err != nil {
		return err
	}
	now := time.Now()
	purchase := Purchase{
		ID:          nextID(purchases),
		Description: description,
		Date:        now,
		Amount:      amount,
	}
	purchases = append(purchases, purchase)
	if err := saveTasks(purchases); err != nil {
		return err
	}
	fmt.Printf("Task added successfully (ID: %d)\n", purchase.ID)
	return nil
}

func cmdUpdate(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: update <id> \"new description\"")
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	task, _ := findTask(tasks, id)
	if task == nil {
		return fmt.Errorf("task %d not found", id)
	}
	task.Description = strings.Join(args[2:], " ")

	if err := saveTasks(tasks); err != nil {
		return err
	}
	fmt.Println("Task updated")
	return nil
}

func cmdDelete(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: delete <id>")
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	_, idx := findTask(tasks, id)
	if idx == -1 {
		return fmt.Errorf("task %d not found", id)
	}
	tasks = append(tasks[:idx], tasks[idx+1:]...)
	if err := saveTasks(tasks); err != nil {
		return err
	}
	fmt.Println("Task deleted")
	return nil
}

func cmdSetStatus(args []string, status string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: mark-%s <id>", status)
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	task, _ := findTask(tasks, id)
	if task == nil {
		return fmt.Errorf("task %d not found", id)
	}

	if err := saveTasks(tasks); err != nil {
		return err
	}
	fmt.Printf("Task %d marked as %s\n", id, status)
	return nil
}

//func cmdList(args []string) error {
//	tasks, err := loadTasks()
//	if err != nil {
//		return err
//	}
//	filter := ""
//	if len(args) > 1 {
//		filter = args[1]
//	}
//	for _, t := range tasks {
//		if filter == "" || t.Status == filter {
//			fmt.Printf("[%d] %s [%s] (created: %s, updated: %s)\n", t.ID, t.Description, t.Status, t.CreatedAt.Format("2006-01-02 15:04"), t.UpdatedAt.Format("2006-01-02 15:04"))
//		}
//	}
//	return nil
//}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-cli <command> [arguments]")
		return
	}
	cmd := os.Args[1]
	args := os.Args[1:]
	var err error
	switch cmd {
	case "add":
		err = cmdAdd(args)
	case "update":
		err = cmdUpdate(args)
	case "delete":
		err = cmdDelete(args)
	case "mark-todo":
		err = cmdSetStatus(args, StatusTodo)
	case "mark-in-progress":
		err = cmdSetStatus(args, StatusInProgress)
	case "mark-done":
		err = cmdSetStatus(args, StatusDone)
	//case "list":
	//	err = cmdList(args)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
	}
	if err != nil {
		fmt.Println("Error:", err)
	}
}
