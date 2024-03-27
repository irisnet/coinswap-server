package task

import (
	"fmt"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/robfig/cron/v3"
	"time"
)

type Task interface {
	Name() string
	Cron() int // second of Intervals
	Start()
	DoTask(fn func(string) chan bool) error
}

var (
	tasks []Task
)

func RegisterTasks(task ...Task) {
	tasks = append(tasks, task...)
}

// GetTasks get all the task
func GetTasks() []Task {
	return tasks
}

func Start() {

	if len(GetTasks()) == 0 {
		return
	}

	for _, one := range GetTasks() {
		var taskId = fmt.Sprintf("%s[%s]", one.Name(), FmtTime(time.Now(), DateFmtYYYYMMDD))
		logger.Info("timerTask begin to work", logger.String("taskId", taskId))
		go one.Start()
	}

	// tasks manager by cron job
	c := cron.New()
	// add cronjob
	c.Start()

}
