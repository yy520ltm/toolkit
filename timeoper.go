package toolkit

import (
	"time"
)

//获取当前时间直接之前或者之后 before，after
//之前 -10s ，之后  10s
func GetTimeBOrA(nowTime time.Time, timeStr string) time.Time {
	t1, _ := time.ParseDuration(timeStr)
	return nowTime.Add(t1)
}

//两个时间相差多少时间单位
func GetTimeSub(t1, t2 time.Time) time.Duration {
	return t2.Sub(t1)
}

//根据时间值返回为float64值
func GetTimeNum(timeDuration,timeDuration2 time.Duration) float64 {

	switch timeDuration {
	case time.Hour:
		return timeDuration2.Hours()
	case time.Minute:
		return timeDuration2.Minutes()
	case time.Second:
		return timeDuration2.Seconds()
		//不写了，要写你自己写
	}
	return 0
}
