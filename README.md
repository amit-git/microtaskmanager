# microtaskmanager
A Lightweight command line tool for managing daily task list.

> v <time range> : shows tasks filtered by time range
>> a "task description" - adds a new task
>> r {taskId} - removes task
>> d {taskId} - marks the task done
>> b {taskId} - puts the task in a backlog
>> w {taskId} - start work on the task, moves it out of backlog if needed

> v now-2 : shows tasks finished in last 2 days
> v backlog : shows tasks in backlog
