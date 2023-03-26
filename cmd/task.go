package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type task struct {
	name           string
	spinner        spinner.Spinner
	finalMSGFormat string
}

type Tasks struct {
	tasks   []*task
	current int
}

func NewTasks(names []string, opts ...spinner.Option) Tasks {
	tasks := make([]*task, 0, len(names))
	for i, name := range names {
		t := &task{name: name, spinner: *spinner.New(spinner.CharSets[14], 100*time.Millisecond)}
		for _, opt := range opts {
			opt(&t.spinner)
		}
		t.spinner.Prefix = " "
		t.spinner.Suffix = fmt.Sprintf(" [%d/%d] %s...", i+1, len(names), name)
		t.finalMSGFormat = fmt.Sprintf("   [%d/%d] %s (%%s)\n", i+1, len(names), name)
		tasks = append(tasks, t)
	}
	return Tasks{tasks: tasks}
}

// Start starts the first task
func (t *Tasks) Start() {
	t.current = 0
	t.tasks[t.current].spinner.Start()
}

// Next stops the current task and starts the next task
func (t *Tasks) Next() {
	t.tasks[t.current].spinner.FinalMSG = fmt.Sprintf(t.tasks[t.current].finalMSGFormat, "Done")
	t.tasks[t.current].spinner.Stop()
	t.current++
	if t.current < len(t.tasks) {
		t.tasks[t.current].spinner.Start()
	}
}

// Error stops the current task
func (t *Tasks) Error() {
	t.tasks[t.current].spinner.FinalMSG = fmt.Sprintf(t.tasks[t.current].finalMSGFormat, "Error")
	t.tasks[t.current].spinner.Stop()
}

func (t *Tasks) Skip() {
	t.tasks[t.current].spinner.FinalMSG = fmt.Sprintf(t.tasks[t.current].finalMSGFormat, "Skipped")
	t.tasks[t.current].spinner.Stop()
	t.current++
	if t.current < len(t.tasks) {
		t.tasks[t.current].spinner.Start()
	}
}

// Close stops all active tasks
func (t *Tasks) Close() {
	for _, task := range t.tasks {
		if task.spinner.Active() {
			t.tasks[t.current].spinner.FinalMSG = fmt.Sprintf(t.tasks[t.current].finalMSGFormat, "Error")
			task.spinner.Stop()
		}
	}
}
