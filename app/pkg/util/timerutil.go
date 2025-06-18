package timerutil

import (
	"log"
	"sync"
	"time"
)

type TimerUtil struct {
	timers sync.Map // map[string]time.Time
}

func NewTimerUtil() *TimerUtil {
	return &TimerUtil{}
}

var D = NewTimerUtil()

const defaultFlag = "default"

// StartTimer 默认标识.
func (t *TimerUtil) StartTimer() {
	t.StartTimerWithFlag(defaultFlag)
}

// StartTimerWithFlag 带标识启动计时.
func (t *TimerUtil) StartTimerWithFlag(flag string) {
	t.timers.Store(flag, time.Now())
}

// GetTimer 返回从启动到当前的耗时，单位毫秒.
func (t *TimerUtil) GetTimer() int64 {
	return t.GetTimerWithFlag(defaultFlag)
}

func (t *TimerUtil) GetTimerWithFlag(flag string) int64 {
	if start, ok := t.timers.Load(flag); ok {
		if startTime, ok2 := start.(time.Time); ok2 {
			return time.Since(startTime).Milliseconds()
		}
	}
	return -1 // 未找到计时，返回-1表示错误
}

// StopTimer 打印所有计时信息.
func (t *TimerUtil) StopTimer() {
	stop := time.Now()
	log.Println("Stop timer.")
	t.timers.Range(func(key, value interface{}) bool {
		flag, ok1 := key.(string)
		startTime, ok2 := value.(time.Time)
		if ok1 && ok2 {
			elapsed := stop.Sub(startTime).Milliseconds()
			log.Printf("flag = %s, spend time %d ms.\n", flag, elapsed)
		}
		return true
	})
}
