package main

import (
	"strings"
	"time"
)

func doView(tasks TaskList) {
	order, reversed := OrderFromString(*orderFlag)
	options := &ViewOptions{
		ShowAll:   *allFlag,
		Summarise: *summaryFlag,
		Order:     order,
		Reversed:  reversed,
	}
	view := NewConsoleView()
	view.ShowTree(tasks, options)
}

func doAdd(tasks TaskList, graft TaskNode, priority Priority, text string) {
	graft.Create(text, priority)
	saveTaskList(tasks)
}

func doMarkDone(tasks TaskList, references []Task) {
	for _, task := range references {
		task.SetCompleted()
	}
	saveTaskList(tasks)
}

func doMarkNotDone(tasks TaskList, references []Task) {
	for _, task := range references {
		task.SetCompletionTime(time.Time{})
	}
	saveTaskList(tasks)
}

func doReparent(tasks TaskList, task TaskNode, below TaskNode) {
	ReparentTask(task, below)
	saveTaskList(tasks)
}

func doRemove(tasks TaskList, references []Task) {
	for _, task := range references {
		task.Delete()
	}
	saveTaskList(tasks)
}

func doPurge(tasks TaskList, age time.Duration) {
	cutoff := time.Now().Add(-age)
	matches := tasks.FindAll(func(task Task) bool {
		return !task.CompletionTime().IsZero() && task.CompletionTime().Before(cutoff)
	})
	for _, m := range matches {
		m.Delete()
	}
	saveTaskList(tasks)
}

func doSetTitle(tasks TaskList, args []string) {
	title := strings.Join(args, " ")
	tasks.SetTitle(title)
	saveTaskList(tasks)
}

func doShowInfo(tasks TaskList, index string) {
	task := tasks.Find(index)
	if task == nil {
		fatalf("no such task %s", index)
	}
	view := NewConsoleView()
	view.ShowTaskInfo(task)
}

func doEditTask(tasks TaskList, task Task, priority Priority, text string) {
	if text != "" {
		task.SetText(text)
	}
	if priority != -1 {
		task.SetPriority(priority)
	}
	saveTaskList(tasks)
}

func editTask(tasks TaskList, priority Priority) {
	if len(*taskText) < 1 {
		fatalf("expected [-p <priority>] <task> [<text>]")
	}
	task := tasks.Find((*taskText)[0])
	if task == nil {
		fatalf("invalid task %s", (*taskText)[0])
	}
	text := strings.Join((*taskText)[1:], " ")
	if *priorityFlag == "" {
		priority = -1
	}
	doEditTask(tasks, task, priority, text)
}
