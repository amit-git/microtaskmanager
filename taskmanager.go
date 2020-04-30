package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TasksContext struct {
	currentTasksFile  *os.File
	finishedTasksFile *os.File
	backlogFile       *os.File
}

type Task struct {
	id          int
	description string
}

const TimeFormat = "01-02-2006"

// global state
var currentTaskId = 0
var tasks = map[int]*Task{}
var backlogTasks = make([]string, 0)
var ctx *TasksContext
var currentFilePathFull string
var backlogFilePathFull string

//var doneTasks map[*time.Time][]string

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

	currentFilePathFull = tasksFolder + "/current"
	backlogFilePathFull = tasksFolder + "/backlog"

	ctx = &TasksContext{
		currentTasksFile:  getFile(currentFilePathFull),
		finishedTasksFile: getFile(tasksFolder + "/done"),
		backlogFile:       getFile(backlogFilePathFull),
	}

	loadCurrentTasks()

	for {
		processOption(presentOptions())
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

func processOption(cmd string) {
	if cmd == "q" {
		fmt.Println("Ok. Bye")
		_ = ctx.currentTasksFile.Close()
		_ = ctx.finishedTasksFile.Close()
		_ = ctx.backlogFile.Close()
		os.Exit(0)
	}

	if strings.HasPrefix(cmd, "a ") {
		addTask(strings.Replace(cmd, "a ", "", 1))
		return
	}
	if strings.HasPrefix(cmd, "v") {
		viewCurrentTasks()
		return
	}
	if strings.HasPrefix(cmd, "d") {
		taskIdStr := strings.Replace(cmd, "d", "", 1)
		taskId, err := strconv.Atoi(taskIdStr)
		if err != nil {
			fmt.Printf("Invalid taskId %s :: %s", cmd, err.Error())
			return
		}
		if tasks[taskId] == nil {
			fmt.Println("Invalid taskId " + string(taskId))
			return
		}
		markTaskDone(taskId)
		return
	}

	fmt.Println("Default Processing :: " + cmd)
}

func addTask(taskDescription string) {
	currentTaskId = currentTaskId + 1
	tasks[currentTaskId] = &Task{
		id:          currentTaskId,
		description: taskDescription,
	}
	saveCurrentTasks()
	fmt.Printf("New Task with id %d saved\n", currentTaskId)
}

func markTaskDone(tid int) {
	currentTime := time.Now()
	recordedTaskLine := fmt.Sprintf("%s\t%s\n", tasks[tid].description, currentTime.Format(TimeFormat))
	if _, err := ctx.finishedTasksFile.WriteString(recordedTaskLine); err != nil {
		fmt.Printf("Unable to write to %s :: %s\n", ctx.finishedTasksFile.Name(), err.Error())
		os.Exit(-1)
	}
	fmt.Printf("Task '%s' DONE", tasks[tid].description)
	delete(tasks, tid)
	saveCurrentTasks()
}

func viewCurrentTasks() {
	for _, tid := range getSortedTaskIds() {
		fmt.Printf("%d. %s\n", tasks[tid].id, tasks[tid].description)
	}
}

func getSortedTaskIds() []int {
	taskIds := make([]int, 0, len(tasks))
	for k := range tasks {
		taskIds = append(taskIds, k)
	}

	sort.Ints(taskIds)
	return taskIds
}

func getFile(fileName string) *os.File {
	filePtr, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0664)
	if err != nil {
		fmt.Printf("Unable to open %s :: %s", fileName, err.Error())
		os.Exit(-1)
	}
	return filePtr
}

func loadCurrentTasks() {
	fileScanner := bufio.NewScanner(ctx.currentTasksFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		currentTaskId = currentTaskId + 1
		tasks[currentTaskId] = &Task{
			id:          currentTaskId,
			description: fileScanner.Text(),
		}
	}
	fmt.Println("Current tasks list loaded.")
}

func loadBacklog() {
	fileScanner := bufio.NewScanner(ctx.backlogFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		backlogTasks = append(backlogTasks, fileScanner.Text())
	}
	fmt.Println("Backlog file loaded.")
}

func saveCurrentTasks() {
	var lines string
	for _, t := range tasks {
		lines = lines + fmt.Sprintf("%s\n", t.description)
	}

	if err := ioutil.WriteFile(currentFilePathFull, []byte(lines), 0644); err != nil {
		fmt.Println("Error writing file system for current tasks file " + err.Error())
	}
}

func saveBacklog() {
	var lines string
	for _,t := range backlogTasks {
		lines = lines + fmt.Sprintf("%s\n", t)
	}
	if err := ioutil.WriteFile(backlogFilePathFull, []byte(lines), 0644); err != nil {
		fmt.Println("Error writing file system for backlog file " + err.Error())
	}
}
