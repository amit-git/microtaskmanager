package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type TasksContext struct {
	currentTasksFile  *os.File
	finishedTasksFile *os.File
	backlogFile       *os.File
}

type Task struct {
	id          uint8
	description string
}

const TimeFormat = "01-02-2006"

var currentTaskId uint8 = 0
var tasks []*Task

func main() {
	currentTime := time.Now()
	fmt.Println("Showing tasks for ", currentTime.Format(TimeFormat))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Can't find current user information. " + err.Error())
		os.Exit(-1)
	}

	tasksFolder := homeDir + "/tm"
	if _, err := os.Stat(tasksFolder); os.IsNotExist(err) {
		if err := os.Mkdir(tasksFolder, os.ModePerm); err != nil {
			fmt.Println("Error creating " + tasksFolder + " :: " + err.Error())
			os.Exit(-1)
		}
	}

	tasksContext := TasksContext{
		currentTasksFile:  getFile(tasksFolder + "/current"),
		finishedTasksFile: getFile(tasksFolder + "/done"),
		backlogFile:       getFile(tasksFolder + "/backlog"),
	}

	loadCurrentTasks(&tasksContext)

	for {
		processOption(presentOptions(), &tasksContext)
	}
}

func presentOptions() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("")
	fmt.Println("Pick one.")
	fmt.Println("v : view current tasks")
	fmt.Println("a <txt> : adds a task ")
	fmt.Println("b <task-id> : puts a task in the backlog")
	fmt.Println("w <task-id> : moves a task out of backlog into current list")
	fmt.Println("d <task-id> : marks task DONE")
	fmt.Println("q : Quit")
	fmt.Printf("\n> ")
	cmd, e := reader.ReadString('\n')
	if e != nil {
		fmt.Println("Error in reading command " + e.Error())
		os.Exit(-1)
	}
	return strings.Replace(cmd, "\n", "", -1)
}

func processOption(cmd string, ctx *TasksContext) {
	if cmd == "q" {
		fmt.Println("Ok. Bye")
		os.Exit(0)
	}

	if strings.HasPrefix(cmd, "a ") {
		addTask(strings.Replace(cmd, "a ", "", 1), ctx)
		return
	}
	if strings.HasPrefix(cmd, "v") {
		viewCurrentTasks()
		return
	}

	fmt.Println("Default Processing :: " + cmd)
}

func addTask(taskDescription string, ctx *TasksContext) {
	if _, err := ctx.currentTasksFile.WriteString(taskDescription + "\n"); err != nil {
		fmt.Printf("Unable to write to %s :: %s\n", ctx.currentTasksFile.Name(), err.Error())
		os.Exit(-1)
	}

	currentTaskId = currentTaskId + 1
	tasks = append(tasks, &Task{
		id:          currentTaskId,
		description: taskDescription,
	})
	fmt.Println("New Task saved")
}

func viewCurrentTasks() {
	for _, t := range tasks {
		fmt.Printf("%d. %s\n", t.id, t.description)
	}
}

func getFile(fileName string) *os.File {
	filePtr, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		fmt.Printf("Unable to open %s :: %s", fileName, err.Error())
		os.Exit(-1)
	}
	return filePtr
}

func loadCurrentTasks(ctx *TasksContext) {
	fileScanner := bufio.NewScanner(ctx.currentTasksFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		currentTaskId = currentTaskId + 1
		tasks = append(tasks, &Task{
			id:          currentTaskId,
			description: fileScanner.Text(),
		})
	}
	fmt.Println("Current tasks list loaded.")
}
