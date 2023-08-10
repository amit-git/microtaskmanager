# microtaskmanager
A Lightweight command line tool for managing daily task list.

## Three modes of operation

1. Current / TODO tasks mode for viewing and editing current list of tasks
1. Backlog mode for viewing and editing backlog items
1. Finished tasks viewing mode

Depending on each mode, the following commands are available

### Current / Active task list

1. vc : shows current task list
1. vb : shows backlog
1. vd now-<X>d : shows tasks in last X days
1. a "task description" - adds a new task
1. r {taskId} - removes task
1. d {taskId} - marks the task done
1. e {taskId} - edits the current task (Coming soon)
1. b {taskId} - puts the task in a backlog
1. q : quit
1. h : show commands available in this mode

### Backlog 

1. vc : shows current task list
1. vb : shows backlog
1. vd now-<X>d : shows tasks in last X days
1. a "task description" - adds a new task
1. r {taskId} - removes task
1. w {taskId} - start work on the task, moves it out of backlog if needed
1. q : quit
1. h : show commands available in this mode

### Finished task list view

1. vc : shows current task list
1. vb : shows backlog
1. vd now-<X>d : shows tasks in last X days
1. q : quit
1. h : show commands available in this mode
