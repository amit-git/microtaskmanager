# microtaskmanager
A Lightweight command line tool for managing daily task list.

## Three modes of operation

1. Current / TODO tasks mode for viewing and editing current list of tasks
1. Backlog mode for viewing and editing backlog items
1. Finished tasks viewing mode

Depending on each mode, the following commands are available

### Current / Active task list

0. vc : shows current task list
0. vb : shows backlog
0. vd now-<X>d : shows tasks in last X days
0. a "task description" - adds a new task
0. r {taskId} - removes task
0. d {taskId} - marks the task done
0. e {taskId} - edits the current task (Not implemented yet)
0. b {taskId} - puts the task in a backlog
0. q : quit
0. h : show commands available in this mode

### Backlog 

0. vc : shows current task list
0. vb : shows backlog
0. vd now-<X>d : shows tasks in last X days
0. a "task description" - adds a new task
0. r {taskId} - removes task
0. w {taskId} - start work on the task, moves it out of backlog if needed
0. q : quit
0. h : show commands available in this mode

### Finished task list view

0. vc : shows current task list
0. vb : shows backlog
0. vd now-<X>d : shows tasks in last X days
0. q : quit
0. h : show commands available in this mode
