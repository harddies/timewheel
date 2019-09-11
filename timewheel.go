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

func (tp *taskPool) Schedule(name string, d ITimeWheelData, tf TriggerFunc, tm time.Time, isClearOld bool) {
	var (
		eq, tok = tp.queue[name]
		now     = time.Now()
	)

	if eq.Equal(tm) {
		return
	}

	if !tm.After(now) {
		if isClearOld && tok {
			tp.stopCh[name] <- true
			delete(tp.queue, name)
			delete(tp.stopCh, name)
		}
		tf(d)
		return
	}

	if tok {
		tp.stopCh[name] <- true
	}
	tp.queue[name] = tm
	tp.stopCh[name] = trigger(name, d, tf, tm)
	fmt.Printf("TimeWheel will excute tf(%p) at %+v\n", tf, tm)
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
