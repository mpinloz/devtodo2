package main

import (
	"strings"
	"testing"
	"time"
)

func getTestTaskList() TaskList {
	tasks := NewTaskList()
	tasks.Create("do A", MEDIUM)
	tasks.Create("do B", MEDIUM)
	tasks.Create("do C", MEDIUM)
	return tasks
}

func getReferences(tasks TaskList, index string) []Task {
	textArray := strings.Split(index, " ")
	taskText = &textArray
	references := resolveTaskReferences(tasks, *taskText)
	return references
}

func TestDoAdd(t *testing.T) {
	tasks := getTestTaskList()
	graft := tasks.Find("1")
	priority := PriorityFromString("low")
	text := "do A.A"
	doAdd(tasks, graft, priority, text)
	if tasks == nil || tasks.Find("1.1").Text() != "do A.A" {
		t.Fail()
	}
}

func TestDoMarkDone(t *testing.T) {
	tasks := getTestTaskList()
	references := getReferences(tasks, "2")
	doMarkDone(tasks, references)
	if tasks == nil || tasks.At(1).CompletionTime().IsZero() {
		t.Fail()
	}
}

func TestDoMarkNotDone(t *testing.T) {
	tasks := getTestTaskList()
	references := getReferences(tasks, "2")
	doMarkDone(tasks, references)
	doMarkNotDone(tasks, references)
	if tasks == nil || !tasks.At(1).CompletionTime().IsZero() {
		t.Fail()
	}
}

func TestDoReparent(t *testing.T) {
	tasks := getTestTaskList()
	childTask := tasks.At(2)
	parentTask := tasks.At(0)
	doReparent(tasks, childTask, parentTask)
	if tasks == nil || tasks.Find("1.1") == nil || tasks.Find("1.1").Text() != "do C" {
		t.Fail()
	}
}

func TestDoSetTitle(t *testing.T) {
	title := "Titulo test"
	tasks := NewTaskList()
	doSetTitle(tasks, strings.Split(title, " "))
	if tasks == nil || tasks.Title() != title {
		t.Fail()
	}
}

func TestDoRemove(t *testing.T) {
	tasks := getTestTaskList()
	text := "1-2"
	textArray := strings.Split(text, " ")
	taskText = &textArray
	references := resolveTaskReferences(tasks, *taskText)
	doRemove(tasks, references)
	if tasks == nil || tasks.Len() != 1 || tasks.At(0).Text() != "do C" {
		t.Fail()
	}
}

func TestDoPurge(t *testing.T) {
	tasks := getTestTaskList()
	references := getReferences(tasks, "2")
	doMarkDone(tasks, references)
	age, err := time.ParseDuration("-2s")
	if err != nil {
		fatalf("error")
	}
	doPurge(tasks, age)
	if tasks == nil || tasks.Len() != 2 {
		t.Fail()
	}
}

func TestDoEditTask(t *testing.T) {
	tasks := getTestTaskList()
	text := "do test"
	priority := PriorityFromString("high")
	task := tasks.At(0)
	doEditTask(tasks, task, priority, text)
	if tasks == nil || tasks.At(0).Text() != "do test" || tasks.At(0).Priority().String() != "high" {
		t.Fail()
	}
}

func TestEditTask(t *testing.T) {
	tasks := getTestTaskList()
	text := strings.Split("1 do test", " ")
	priority := PriorityFromString("high")
	editTask(tasks, priority, text)
	if tasks == nil || tasks.At(0).Text() != "do test" || tasks.At(0).Priority().String() != "high" {
		t.Fail()
	}

	text = strings.Split("1 do test without priority", " ")
	editTask(tasks, -1, text)
	if tasks == nil || tasks.At(0).Text() != "do test without priority" {
		t.Fail()
	}
}

func TestAddTask(t *testing.T) {
	tasks := getTestTaskList()
	graft := "1"
	priority := PriorityFromString("low")
	taskText := strings.Split("do A.A", " ")
	addTask(tasks, taskText, priority, graft)
	if tasks == nil || tasks.Find("1.1").Text() != "do A.A" {
		t.Fail()
	}
}

func TestReparenTask(t *testing.T) {
	tasks := getTestTaskList()
	reparentTask(tasks, strings.Split("3 1", " "))
	if tasks == nil || tasks.Find("1.1") == nil || tasks.Find("1.1").Text() != "do C" {
		t.Fail()
	}

	reparentTask(tasks, strings.Split("1.1", " "))
	if tasks == nil || tasks.Find("3") == nil || tasks.Find("3").Text() != "do C" {
		t.Fail()
	}
}
