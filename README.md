#Tasky

A CLI TODO app built with pure Go

Usage:
```
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
	tasky done <task id> (Assigns the status 'done' to the task with the specified id)
```
