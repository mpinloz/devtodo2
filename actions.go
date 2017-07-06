package main

import (
	"flag"
	"fmt"
	"os"
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
}

func doMarkDone(tasks TaskList, references []Task) {
	for _, task := range references {
		task.SetCompleted()
	}
}

func doMarkNotDone(tasks TaskList, references []Task) {
	for _, task := range references {
		task.SetCompletionTime(time.Time{})
	}
}

func doReparent(tasks TaskList, task TaskNode, below TaskNode) {
	ReparentTask(task, below)
}

func doRemove(tasks TaskList, references []Task) {
	for _, task := range references {
		task.Delete()
	}
}

func doPurge(tasks TaskList, age time.Duration) {
	cutoff := time.Now().Add(-age)
	matches := tasks.FindAll(func(task Task) bool {
		return !task.CompletionTime().IsZero() && task.CompletionTime().Before(cutoff)
	})
	for _, m := range matches {
		m.Delete()
	}
}

func doSetTitle(tasks TaskList, args []string) {
	title := strings.Join(args, " ")
	tasks.SetTitle(title)
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
}

func editTask(tasks TaskList, priority Priority, taskText []string) {
	if len(taskText) < 1 {
		fatalf("expected [-p <newpriority>] <task> [<newtext>]")
	}
	task := tasks.Find((taskText)[0])
	if task == nil {
		fatalf("invalid task %s", (taskText)[0])
	}
	text := strings.Join((taskText)[1:], " ")
	doEditTask(tasks, task, priority, text)
}

func reparentTask(tasks TaskList, taskText []string) {
	if len(taskText) < 1 {
		fatalf("expected <task> [<new-parent>] for reparenting")
	}
	var below TaskNode
	if len(taskText) == 2 {
		below = resolveTaskReference(tasks, (taskText)[1])
	} else {
		below = tasks
	}
	doReparent(tasks, resolveTaskReference(tasks, (taskText)[0]), below)
}

func addTask(tasks TaskList, taskText []string, priority Priority, graftFlag string) {
	var graft TaskNode = tasks
	if graftFlag != "root" {
		if graft = tasks.Find(graftFlag); graft == nil {
			fatalf("invalid graft index '%s'", graftFlag)
		}
	}
	if len(taskText) == 0 {
		fatalf("expected text for new task")
	}
	text := strings.Join(taskText, " ")
	doAdd(tasks, graft, priority, text)
}

func doManPage(usage string) {
	title := "todo2 1 \"2.2.0\""
	name := "todo2 - Terminal task manager"
	synopsis := "todo2 [options] [args..]"
	author := "Alec Thomas <alec@swapoff.org>"
	fmt.Fprintln(os.Stdout, ".TH "+title)
	fmt.Fprintln(os.Stdout, ".SH NAME\n"+name)
	fmt.Fprintln(os.Stdout, ".SH SYNOPSIS\n"+synopsis)
	fmt.Fprintln(os.Stdout, ".SH DESCRIPTION\n"+usage)
	fmt.Fprintln(os.Stdout, ".SH OPTIONS")

	flag.VisitAll(func(arg1 *flag.Flag) {
		fmt.Fprintf(os.Stdout, ".TP\n.BR \\-%s \n%s\n", arg1.Name, arg1.Usage)
	})
	fmt.Fprintln(os.Stdout, ".SH AUTHOR\n"+author)

}
