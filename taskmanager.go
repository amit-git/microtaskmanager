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

type Mode int

const (
	Current = iota
	Backlog
	Done
)

func (m Mode) String() string {
	return []string{"Current", "Backlog", "Done"}[m]
}

type TasksContext struct {
	currentTasksFile *os.File
	doneTasksFile    *os.File
	backlogFile      *os.File
}

type Task struct {
	id          int
	description string
}

type DoneTask struct {
	description string
	finishedOn  time.Time
}

const TimeFormat = "01-02-2006"

// global state
var mode Mode = Current
var curLatestTaskId = 0
var backlogLatestTaskId = 0
var tasks = map[int]*Task{}
var backlogTasks = map[int]*Task{}
var doneTasks = make([]*DoneTask, 0)
var ctx *TasksContext
var currentFilePathFull string
var backlogFilePathFull string
var doneFilePathFull string

func main() {
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
	doneFilePathFull = tasksFolder + "/done"

	ctx = &TasksContext{
		currentTasksFile: getFile(currentFilePathFull),
		doneTasksFile:    getFile(doneFilePathFull),
		backlogFile:      getFile(backlogFilePathFull),
	}

	// load persisted data
	loadCurrentTasks()
	loadBacklog()
	loadDoneTasks()

	// first view is current list of tasks
	viewCurrentTasks()
	for {
		cmd := readCmd()
		processOption(cmd)
	}
}

func readCmd() string {
	fmt.Printf("\n> ")
	reader := bufio.NewReader(os.Stdin)
	cmd, e := reader.ReadString('\n')
	if e != nil {
		fmt.Println("Error in reading command " + e.Error())
		os.Exit(-1)
	}
	return strings.Replace(cmd, "\n", "", -1)
}

func showHelp() string {
	switch mode {
	case Current:
		showCurrentModeOptions()
	case Backlog:
		showBacklogModeOptions()
	case Done:
		showDoneModeOptions()
	default:
		fmt.Println("Invalid Mode " + mode.String())
		os.Exit(-1)
	}
	return ""
}

func showCurrentModeOptions() {
	fmt.Println("")
	fmt.Println("vc : view current tasks")
	fmt.Println("vb : view backlog tasks")
	fmt.Println("vd now-Xd : view done tasks in past X days")
	fmt.Println("a <txt> : adds a task ")
	fmt.Println("r <task-id> : removes a task from current list")
	fmt.Println("b <task-id> : puts a task in the backlog")
	fmt.Println("d <task-id> : marks task DONE")
	fmt.Println("q : Quit")
}

func showBacklogModeOptions() {
	fmt.Println("")
	fmt.Println("vc : view current tasks")
	fmt.Println("vb : view backlog tasks")
	fmt.Println("vd now-Xd : view done tasks in past X days")
	fmt.Println("a <txt> : adds a task to backlog")
	fmt.Println("w <task-id> : moves a task out of backlog into current list")
	fmt.Println("r <task-id> : removes a task from the backlog")
	fmt.Println("q : Quit")
}

func showDoneModeOptions() {
	fmt.Println("")
	fmt.Println("vb : view current backlog")
	fmt.Println("vc : view current tasks")
	fmt.Println("vd now-Xd : view done tasks in past X days")
	fmt.Println("q : Quit")
}

func processOption(cmd string) {
	if cmd == "h" {
		showHelp()
		return
	}

	if cmd == "q" {
		fmt.Println("Ok. Bye")
		_ = ctx.currentTasksFile.Close()
		_ = ctx.doneTasksFile.Close()
		_ = ctx.backlogFile.Close()
		os.Exit(0)
	}

	// View Commands
	if strings.HasPrefix(cmd, "vc") {
		viewCurrentTasks()
		return
	}
	if strings.HasPrefix(cmd, "vb") {
		viewBacklog()
		return
	}

	if strings.HasPrefix(cmd, "vd") {
		nowStr := strings.Replace(cmd, "vd ", "", 1)
		if !strings.HasPrefix(nowStr, "now-") {
			fmt.Println("Invalid time expression " + nowStr)
			return
		}
		durationStr := strings.Replace(nowStr, "now-", "", 1)
		if !strings.HasSuffix(durationStr, "d") {
			fmt.Println("Invalid duration expression " + durationStr)
			return
		}
		daysStr := strings.Replace(durationStr, "d", "", 1)
		daysNum, err := strconv.Atoi(daysStr)
		if err != nil {
			fmt.Println("Invalid number format " + daysStr)
			return
		}
		duration := time.Duration(-daysNum) * 24 * time.Hour
		viewDoneTasks(duration)
	}

	// Edit commands
	if strings.HasPrefix(cmd, "a ") {
		addTask(strings.Replace(cmd, "a ", "", 1))
		return
	}

	if strings.HasPrefix(cmd, "r") {
		taskIdStr := strings.Replace(cmd, "r", "", 1)
		taskId, err := strconv.Atoi(taskIdStr)
		if err != nil {
			fmt.Printf("Invalid taskId %s :: %s", cmd, err.Error())
			return
		}
		removeTask(taskId)
		return
	}

	// Move commands
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
		viewCurrentTasks()
		return
	}

	if strings.HasPrefix(cmd, "b") {
		taskIdStr := strings.Replace(cmd, "b", "", 1)
		taskId, err := strconv.Atoi(taskIdStr)
		if err != nil {
			fmt.Printf("Invalid taskId %s :: %s", cmd, err.Error())
			return
		}
		if tasks[taskId] == nil {
			fmt.Println("Invalid taskId " + string(taskId))
			return
		}
		moveTaskToBacklog(taskId)
		viewCurrentTasks()
		return
	}

	if strings.HasPrefix(cmd, "w") {
		taskIdStr := strings.Replace(cmd, "w", "", 1)
		taskId, err := strconv.Atoi(taskIdStr)
		if err != nil {
			fmt.Printf("Invalid taskId %s :: %s", cmd, err.Error())
			return
		}
		moveTaskToCurrent(taskId)
		viewBacklog()
		return
	}
}

func moveTaskToCurrent(tid int) {
	if backlogTasks[tid] != nil {
		taskDescription := backlogTasks[tid].description
		addTaskToCurrent(taskDescription)
		saveCurrentTasks()
		removeTaskFromBacklog(tid)
		fmt.Printf("Task '%s' moved to current.\n", taskDescription)
	} else {
		fmt.Printf("Invalid taskId for backlog %d\n", tid)
	}
}

func moveTaskToBacklog(tid int) {
	backlogLatestTaskId = backlogLatestTaskId + 1
	backlogTasks[backlogLatestTaskId] = &Task{
		id:          backlogLatestTaskId,
		description: tasks[tid].description,
	}
	delete(tasks, tid)
	saveCurrentTasks()
	saveBacklog()
	fmt.Printf("Task '%s' moved to backlog.\n", backlogTasks[backlogLatestTaskId].description)
}

func addTask(taskDescription string) {
	switch mode {
	case Current:
		addTaskToCurrent(taskDescription)
		viewCurrentTasks()
	case Backlog:
		addTaskToBacklog(taskDescription)
		viewBacklog()
	case Done:
		fmt.Printf("Invalid state for adding new tasks \n")
	}
}

func removeTask(tid int) {
	switch mode {
	case Current:
		removeTaskFromCurrent(tid)
		viewCurrentTasks()
	case Backlog:
		removeTaskFromBacklog(tid)
		viewBacklog()
	case Done:
		fmt.Printf("Invalid state for removing tasks \n")
	}

}

func removeTaskFromBacklog(tid int) {
	delete(backlogTasks, tid)
	saveBacklog()
	fmt.Printf("Removed\n")
}

func removeTaskFromCurrent(tid int) {
	delete(tasks, tid)
	saveCurrentTasks()
	fmt.Printf("Removed\n")
}

func addTaskToCurrent(taskDescription string) {
	curLatestTaskId = curLatestTaskId + 1
	tasks[curLatestTaskId] = &Task{
		id:          curLatestTaskId,
		description: taskDescription,
	}
	saveCurrentTasks()
	fmt.Printf("Saved\n")
}

func addTaskToBacklog(taskDescription string) {
	backlogLatestTaskId = backlogLatestTaskId + 1
	backlogTasks[backlogLatestTaskId] = &Task{
		id:          backlogLatestTaskId,
		description: taskDescription,
	}
	saveBacklog()
	fmt.Printf("Saved\n")
}

func markTaskDone(tid int) {
	currentTime := time.Now()
	doneTasks = append(doneTasks, &DoneTask{
		description: tasks[tid].description,
		finishedOn:  currentTime,
	})
	fmt.Printf("Task '%s' DONE\n", tasks[tid].description)
	delete(tasks, tid)
	saveCurrentTasks()
	saveDoneTasks()
}

func viewCurrentTasks() {
	mode = Current
	sortedTaskIds := getSortedTaskIds(tasks)
	showTasksHeader()
	for _, tid := range sortedTaskIds {
		fmt.Printf("%d. %s\n", tasks[tid].id, tasks[tid].description)
	}
}

func viewBacklog() {
	mode = Backlog
	sortedTaskIds := getSortedTaskIds(backlogTasks)
	showTasksHeader()
	for _, tid := range sortedTaskIds {
		fmt.Printf("%d. %s\n", backlogTasks[tid].id, backlogTasks[tid].description)
	}
}

func viewDoneTasks(d time.Duration) {
	mode = Done
	threshold := time.Now().Add(d)
	for _, t := range doneTasks {
		if t.finishedOn.After(threshold) {
			fmt.Printf("%s : %s\n", t.finishedOn.Format(TimeFormat), t.description)
		}
	}
}

func getSortedTaskIds(taskList map[int]*Task) []int {
	taskIds := make([]int, 0, len(taskList))
	for k := range taskList {
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
		curLatestTaskId = curLatestTaskId + 1
		tasks[curLatestTaskId] = &Task{
			id:          curLatestTaskId,
			description: fileScanner.Text(),
		}
	}
	fmt.Printf("Loading current tasks :: %d\n", len(tasks))
}

func loadBacklog() {
	fileScanner := bufio.NewScanner(ctx.backlogFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		backlogLatestTaskId = backlogLatestTaskId + 1
		backlogTasks[backlogLatestTaskId] = &Task{
			id:          backlogLatestTaskId,
			description: fileScanner.Text(),
		}
	}
	fmt.Printf("Loading backlog tasks :: %d\n", len(backlogTasks))
}

func loadDoneTasks() {
	fileScanner := bufio.NewScanner(ctx.doneTasksFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lineParts := strings.Split(line, "\t")
		finishedOnTime, err := time.Parse(TimeFormat, lineParts[1])
		if err != nil {
			fmt.Println("Oops, error loading done tasks file " + err.Error())
			os.Exit(-1)
		}
		doneTasks = append(doneTasks, &DoneTask{
			description: lineParts[0],
			finishedOn:  finishedOnTime,
		})
	}
	fmt.Printf("Loading finished tasks :: %d\n", len(doneTasks))
}

func saveCurrentTasks() {
	var lines string
	for _, t := range tasks {
		lines = lines + fmt.Sprintf("%s\n", t.description)
	}

	if err := ioutil.WriteFile(currentFilePathFull, []byte(lines), 0644); err != nil {
		fmt.Println("Error writing file system for current tasks file " + err.Error())
		os.Exit(-1)
	}
}

func saveBacklog() {
	var lines string
	for _, t := range backlogTasks {
		lines = lines + fmt.Sprintf("%s\n", t.description)
	}
	if err := ioutil.WriteFile(backlogFilePathFull, []byte(lines), 0644); err != nil {
		fmt.Println("Error writing file system for backlog file " + err.Error())
		os.Exit(-1)
	}
}

func saveDoneTasks() {
	var lines, recordedTaskLine string
	for _, t := range doneTasks {
		recordedTaskLine = fmt.Sprintf("%s\t%s\n", t.description, t.finishedOn.Format(TimeFormat))
		lines = lines + recordedTaskLine
	}
	if err := ioutil.WriteFile(doneFilePathFull, []byte(lines), 0644); err != nil {
		fmt.Println("Error writing file system for finished tasks file " + err.Error())
		os.Exit(-1)
	}
}

func showTasksHeader() {
	fmt.Printf("\nTasks (%s)----------------------------------------\n", mode.String())
}
