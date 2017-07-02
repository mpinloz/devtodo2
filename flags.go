package main

import (
	"flag"
	"time"
)

// Actions
var addFlag = flag.Bool("a", false, "Add a task.")
var updateFlag = flag.Bool("u", false, "Edit a task, replacing its text.")
var removeFlag = flag.Bool("rm", false, "Remove the given tasks.")
var markDoneFlag = flag.Bool("d", false, "Mark the given tasks as done.")
var markNotDoneFlag = flag.Bool("D", false, "Mark the given tasks as not done.")
var infoFlag = flag.Bool("i", false, "Show information on a task.")
var titleFlag = flag.Bool("t", false, "Set the task list title.")
var reparentFlag = flag.Bool("mv", false, "Reparent task A below task B")
var priorityFlag = flag.String("p", "medium", "priority of newly created tasks (veryhigh,high,medium,low,verylow)")
var allFlag = flag.Bool("A", false, "Show all tasks, even completed ones.")

var purgeFlag = flag.Duration("purge", 0*time.Second, "Purge completed tasks older than this.")
var summaryFlag = flag.Bool("s", false, "Summarise tasks to one line.")
var graftFlag = flag.String("g", "root", "Task to graft new tasks to.")
var importFlag = flag.Bool("import", false, "Import and synchronise TODO items from source code.")
var fileFlag = flag.String("file", ".todo2", "File to load task lists from.")
var legacyFileFlag = flag.String("legacy-file", ".todo", "File to load legacy task lists from.")
var orderFlag = flag.String("order", "priority", "Specify display order of tasks (index,created,completed,text,priority,duration,done)")
var helpManFlag = flag.Bool("help-man", false, "Show help.")
