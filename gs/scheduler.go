package gs

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	// 建议不要重复,但目前暂无实际作用,也可不填
	Name string
	// 由于scheduler轮询频率为1秒,间隔时间的最小单位也请不要小于秒
	Period time.Duration
	// 剩余循环次数,若初始设为负数则将无限循环
	Loop int
	// 此函数内请不要使用panic将错误抛出,直接在内部完成错误处理,否则整个程序将会结束
	Func func(*gin.Context)
	// 上次执行时间
	LastTime time.Time
}

// 这里不能用指针，否则会导致并发问题（在循环里直接用循环变量导致的）
func (task Task) Handle() {
	c := &gin.Context{}
	logger := GetLoggerByGinCtx(c).With().Str("task", task.Name).Logger()

	defer func() {
		if r := recover(); r != nil {
			logger.Error().Any("panic", r).Msg("recovered from panic")
		}
		logger.Info().Msg("handle end")
	}()

	logger.Info().Msg("handle begin")
	if task.Func != nil {
		task.Func(c)
	} else {
		logger.Warn().Msg("has no func")
	}
}

var taskList = make([]Task, 0, 8)

func AddTasks(tasks ...Task) {
	taskList = append(taskList, tasks...)
}

func NewOnceTaskSinceNow(name string, period time.Duration, func_ func(*gin.Context)) *Task {
	return &Task{
		Name:     name,
		Period:   period,
		Loop:     1,
		Func:     func_,
		LastTime: time.Now(),
	}
}

func NewInfiniteTaskSinceNow(name string, period time.Duration, func_ func(*gin.Context)) *Task {
	return &Task{
		Name:     name,
		Period:   period,
		Loop:     -1,
		Func:     func_,
		LastTime: time.Now(),
	}
}

func NewInfiniteTaskImmediately(name string, period time.Duration, func_ func(*gin.Context)) *Task {
	return &Task{
		Name:     name,
		Period:   period,
		Loop:     -1,
		Func:     func_,
		LastTime: time.Now().Add(-period),
	}
}

func NewRepeatTaskSinceNow(name string, period time.Duration, loop int, func_ func(*gin.Context)) *Task {
	return &Task{
		Name:     name,
		Period:   period,
		Loop:     loop,
		Func:     func_,
		LastTime: time.Now(),
	}
}

func NewRepeatTaskImmediately(name string, period time.Duration, loop int, func_ func(*gin.Context)) *Task {
	return &Task{
		Name:     name,
		Period:   period,
		Loop:     loop,
		Func:     func_,
		LastTime: time.Now().Add(-period),
	}
}

func NewDailyTask(name string, hour, minute int, func_ func(*gin.Context)) *Task {
	// BUG: 如果设定时间已经过了，那么会立刻执行而不是等到第二天同时间再执行
	now := time.Now()
	return &Task{
		Name:   name,
		Period: 24 * time.Hour,
		Loop:   -1,
		Func:   func_,
		LastTime: time.Date(
			now.Year(), now.Month(), now.Day()-1,
			hour, minute, 0, 0, now.Location(),
		),
	}
}

func init() {
	go func() {
		tick := time.NewTicker(1 * time.Second)
		for range tick.C {
			neoTaskList := make([]Task, 0, len(taskList)/2)
			for _, task := range taskList {
				if task.LastTime.IsZero() {
					task.LastTime = time.Now()
				} else if time.Since(task.LastTime) >= task.Period {
					if task.Loop != 0 { // 小于零则无限循环,故条件不是大于零
						go task.Handle()
						task.LastTime = time.Now()
					}
					if task.Loop > 0 { // 怕负数溢出什么的,故条件不是不等于零
						task.Loop--
					}
				}
				if task.Loop != 0 { // 小于零则无限循环,故条件不是大于零
					neoTaskList = append(neoTaskList, task)
				}
			}
			taskList = neoTaskList
		}
	}()

	GetLoggerByGinCtx(nil).Info().Msg("scheduler init complete")
}
