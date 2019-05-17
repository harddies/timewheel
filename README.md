##**定点执行任务**

###时间一到，立即执行任务的小工具

####Demo
(```)

	package main
	import (
		"fmt"
		"strconv"
		"time"
		
		"my-go/study/timewheel"
	)
	
	var tw = timewheel.New()
	
	type TestCron struct {
		Who    string
		Do     string
		Whom   string
		Remark string
		Time   time.Time
	}
	
	func (tc *TestCron) OnTrigger() {
		fmt.Printf("定时在（%+v）说一句：%s%s%s，%s\n", tc.Time, tc.Who, tc.Do, tc.Whom, tc.Remark)
	}
	
	func main() {
		now := time.Now()
	
		for i := 2; i <= 2*5; i += 2 {
			tc := &TestCron{
				Who:    "我",
				Do:     "要成为",
				Whom:   "大佬",
			}
			tf := func(d timewheel.ITimeWheelData) {
				d.OnTrigger()
			}
			tc.Time = now.Add(time.Duration(i) * time.Second)
			tw.Schedule("test"+strconv.Itoa(i), tc, tf, tc.Time)
		}
	
		var ch chan bool
		for {
			select {
			case <-ch:
				return
			default:
			}
		}
	}
	
(```)