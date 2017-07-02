/*
  Copyright 2011 Alec Thomas

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const usage = `DevTodo2 - a hierarchical command-line task manager

DevTodo is a program aimed specifically at programmers (but usable by anybody
at the terminal) to aid in day-to-day development.

It maintains a list of items that have yet to be completed, one list for each
project directory. This allows the programmer to track outstanding bugs or
items that need to be completed with very little effort.

Items can be prioritised and are displayed in a hierarchy, so that one item may
depend on another.


  todo2 [-A]
    Display (all) tasks.

  todo2 [-p <priority>][-g <idparent>] -a <text>
    Create a new task.

  todo2 -d <index>
    Mark a task as complete.

  todo2 [-p <priority>] -e <task> [<text>]
    Edit an existing task.
`

var taskText *[]string

func processAction(tasks TaskList) {
	priority := PriorityFromString(*priorityFlag)

	switch {
	case *addFlag:
		addTask(tasks, *taskText, priority, *graftFlag)
	case *markDoneFlag:
		doMarkDone(tasks, resolveTaskReferences(tasks, *taskText))
	case *markNotDoneFlag:
		doMarkNotDone(tasks, resolveTaskReferences(tasks, *taskText))
	case *removeFlag:
		doRemove(tasks, resolveTaskReferences(tasks, *taskText))
	case *reparentFlag:
		reparentTask(tasks, *taskText)
	case *titleFlag:
		doSetTitle(tasks, *taskText)
	case *infoFlag:
		if len(*taskText) < 1 {
			fatalf("expected <task> for info")
		}
		doShowInfo(tasks, (*taskText)[0])
	case *importFlag:
		if len(*taskText) < 1 {
			fatalf("expected list of files to import")
		}
		doImport(tasks, *taskText)
	case *updateFlag:
		editTask(tasks, priority, *taskText)
	case *purgeFlag != 0*time.Second:
		doPurge(tasks, *purgeFlag)
	default:
		doView(tasks)
	}
}

func resolveTaskReference(tasks TaskList, index string) Task {
	task := tasks.Find(index)
	if task == nil {
		fatalf("invalid task index %s", index)
	}
	return task
}

func expandRange(indexRange string) []string {
	// This whole function makes me sad. This kind of manipulation of strings and
	// arrays just should not be this verbose.
	//
	// For constrast, in Python:
	//
	// def expand_range(index):
	//   start_index, end = index.split('-')
	//   start_index, start = start_index.rsplit('.', 1)
	//   for i in range(int(start), int(end) + 1):
	//     yield '%s.%s' % (start_index, str(i))
	ranges := strings.Split(indexRange, "-")
	if len(ranges) != 2 {
		return nil
	}
	startIndex := strings.Split(ranges[0], ".")
	start, err := strconv.Atoi(startIndex[len(startIndex)-1])
	if err != nil {
		return nil
	}
	end, err := strconv.Atoi(ranges[1])
	if err != nil {
		return nil
	}
	rangeIndexes := []string{}
	for i := start; i <= end; i++ {
		index := startIndex[:len(startIndex)-1]
		index = append(index, fmt.Sprintf("%d", i))
		rangeIndexes = append(rangeIndexes, strings.Join(index, "."))
	}
	return rangeIndexes
}

func resolveTaskReferences(tasks TaskList, indices []string) []Task {
	references := make([]Task, 0, len(indices))
	for _, index := range indices {
		if strings.Index(index, "-") == -1 {
			task := resolveTaskReference(tasks, index)
			references = append(references, task)
		} else {
			// Expand ranges. eg. 1.2-5 expands to 1.2 1.3 1.4 1.5
			indexes := expandRange(index)
			if indexes == nil {
				fatalf("invalid task range %s", index)
			}
			for _, rangeIndex := range indexes {
				task := resolveTaskReference(tasks, rangeIndex)
				if task != nil {
					references = append(references, task)
				}
			}
		}
	}
	if len(references) == 0 {
		fatalf("no tasks provided to mark done")
	}
	return references
}

func loadTaskList() (tasks TaskList, err error) {
	// Try loading new-style task file
	if file, err := os.Open(*fileFlag); err == nil {
		defer file.Close()
		loader := NewJSONIO()
		return loader.Deserialize(file)
	}
	// Try loading legacy task file
	if file, err := os.Open(*legacyFileFlag); err == nil {
		defer file.Close()
		loader := NewLegacyIO()
		return loader.Deserialize(file)
	}
	return nil, nil
}

func saveTaskList(tasks TaskList) {
	path := *fileFlag
	previous := path + "~"
	temp := path + "~~"
	var serializeError error
	if file, err := os.Create(temp); err == nil {
		defer func() {
			if err = file.Close(); err != nil {
				os.Remove(temp)
			} else {
				if serializeError != nil {
					return
				}
				if _, err = os.Stat(path); err == nil {
					if err = os.Rename(path, previous); err != nil {
						fatalf("unable to rename %s to %s", path, previous)
					}
				}
				if err = os.Rename(temp, path); err != nil {
					fatalf("unable to rename %s to %s", temp, path)
				}
			}
		}()
		writer := NewJSONIO()
		if serializeError = writer.Serialize(file, tasks); serializeError != nil {
			fatalf(serializeError.Error())
		}
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	taskText = &args
	tasks, err := loadTaskList()
	if err != nil {
		//saveTaskList(tasks)
		//fatalf("Error loadTaskList: %s", err)
		fmt.Println("No file found. Creating one...")
	}
	if tasks == nil {
		fmt.Println("no task")
		tasks = NewTaskList()
	}
	processAction(tasks)
	saveTaskList(tasks)
}
