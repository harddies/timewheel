package timewheel

import (
	"fmt"
	"time"
)

func New() *taskPool {
	tp := &taskPool{
		queue:  make(map[string]time.Time),
		stopCh: make(map[string]chan bool),
	}
	return tp
}

type ITimeWheelData interface {
	OnTrigger()
}

type TriggerFunc func(d ITimeWheelData)

type taskPool struct {
	queue  map[string]time.Time
	stopCh map[string]chan bool
}

func (tp *taskPool) Schedule(name string, d ITimeWheelData, tf TriggerFunc, tm time.Time) {
	_, tok := tp.queue[name]
	now := time.Now()
	if tm.Unix() <= now.Unix() {
		tf(d)
		return
	}
	if !tok {
		tp.queue[name] = tm
		tp.stopCh[name] = trigger(name, d, tf, tm)
		fmt.Printf("TimeWheel will excute tf(%p) at %+v\n", tf, tm)
		return
	}
	tp.stopCh[name] <- true
	tp.queue[name] = tm
	tp.stopCh[name] = trigger(name, d, tf, tm)
}

func trigger(name string, d ITimeWheelData, tf TriggerFunc, tm time.Time) chan bool {
	var (
		deltaT = tm.Unix() - time.Now().Unix()
		ticker = time.NewTicker(time.Duration(deltaT) * time.Second)
		isStop bool
		stopCh = make(chan bool)
	)

	go func() {
		select {
		case <-ticker.C:
			tf(d)
			return
		case isStop = <-stopCh:
			if isStop {
				fmt.Printf("name: %s stop\n", name)
				return
			}
		}
	}()

	return stopCh
}
