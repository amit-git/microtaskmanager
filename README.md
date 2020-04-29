# microtaskmanager
A Lightweight command line tool for managing daily task list.

## Following commands are supported

0. v <time range> : shows tasks filtered by time range
1. a "task description" - adds a new task
2. r {taskId} - removes task
3. d {taskId} - marks the task done
4. b {taskId} - puts the task in a backlog
5. w {taskId} - start work on the task, moves it out of backlog if needed

6. v now-2 : shows tasks finished in last 2 days
7. v backlog : shows tasks in backlog
