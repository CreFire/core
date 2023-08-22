package antnet

import (
	"math"
	"time"
)

const (
	SecondsPerMinute = 60
	SecondsPerHour   = 60 * SecondsPerMinute
	SecondsPerDay    = 24 * SecondsPerHour
	SecondsPerWeek   = 7 * SecondsPerDay
	SecondsPerMonth  = 30 * SecondsPerDay
)

var WeekStart int64 = 1514736000 //修正:不同时区不同

func ParseTime(str string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
}

func Date() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func UnixTime(sec, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

func UnixMs() int64 {
	return time.Now().UnixNano() / 1000000
}

func Now() time.Time {
	return time.Now()
}

// 年
func Year(sec, nsec int64) int {
	return time.Unix(sec, nsec).Year()
}

// 月
func Month(sec, nsec int64) int {
	return int(time.Unix(sec, nsec).Month())
}

// 日
func Day(sec, nsec int64) int {
	return time.Unix(sec, nsec).Day()
}

// 返回指定时间所属月的天数
func GetMonthDay(sec, nsec int64) int {
	year := Year(sec, nsec)
	month := Month(sec, nsec)
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			return 30
		} else {
			return 31
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || ((year % 400) == 0) {
			return 29
		} else {
			return 28
		}
	}
}

// 返回指定时分秒的时间(时分秒的格式为"8:00:00")
func DateStr(sec, nsec int64, hour string) int64 {
	dateStr := time.Unix(sec, nsec).Format("2006-01-02")
	timeStr := Sprintf("%v %v", dateStr, hour)
	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// 周数
func Week(sec, nsec int64) int {
	_, week := time.Unix(sec, nsec).ISOWeek()
	return week
}

func UniqueDay(sec, nsec int64) int32 {
	return int32(Atoi(time.Unix(sec, nsec).Format("20060102")))
}

func UniqueWeek(sec, nsec int64) int32 {
	y, w := time.Unix(sec, nsec).ISOWeek()
	return int32(y*100 + w)
}

func UniqueServerWeek(sec int64) int32 {
	return int32(sec-WeekStart)/SecondsPerWeek + 1
}

// 0-6
func WeekDay(sec, nsec int64) time.Weekday {
	return time.Unix(sec, nsec).Weekday()
}

// 1-7
func WeekDayEx(sec, nsec int64) int32 {
	day := int32(WeekDay(sec, nsec))
	if day == 0 {
		day = 7
	}
	return day
}

func UniqueMonth(sec, nsec int64) int32 {
	atime := time.Unix(sec, nsec)
	return int32(atime.Year()*100 + int(atime.Month()))
}

// 当前时间对应的当天的0点时间戳
func ZeroTime(sec, nsec int64) int64 {
	dateStr := time.Unix(sec, nsec).Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	return t.Unix()
}

func SixTime(sec, nsec int64) int64 {
	return AnyTime(sec, 6*SecondsPerHour)
}

func AnyTime(sec int64, delay int32) int64 {
	t := ZeroTime(sec, 0) + int64(delay)
	if t > sec {
		return t - SecondsPerDay
	} else {
		return t
	}
}

func WeekMondayUnix(sec int64, delay int32) int64 {
	t := AnyTime(sec, delay)
	w := WeekDayEx(t, 0)
	return t - (int64(w)-1)*SecondsPerDay
}

func MonthFirstdayUnix(sec int64, delay int32) int64 {
	t := AnyTime(sec, delay)
	m := time.Unix(t, 0).Day()
	return t - (int64(m)-1)*SecondsPerDay
}

func SixUniqueDay(sec, nsec int64) int32 {
	return AnyUniqueDay(sec, 6*SecondsPerHour)
}

func SixUniqueWeek(sec, nsec int64) int32 {
	return AnyUniqueWeek(sec, 6*SecondsPerHour)
}

func FifteenUniqueWeek(sec, nsec int64) int32 {
	return AnyUniqueWeek(sec, 15*SecondsPerHour+SecondsPerDay)
}

func SixUniqueMonth(sec, nsec int64) int32 {
	return AnyUniqueMonth(sec, 6*SecondsPerHour)
}

func AnyUniqueDay(sec int64, delay int32) int32 {
	return int32(Atoi(time.Unix(sec-int64(delay), 0).Format("20060102")))
}

func AnyUniqueWeek(sec int64, delay int32) int32 {
	y, w := time.Unix(sec-int64(delay), 0).ISOWeek()
	return int32(y*100 + w)
}

func AnyUniqueMonth(sec int64, delay int32) int32 {
	atime := time.Unix(sec-int64(delay), 0)
	return int32(atime.Year()*100 + int(atime.Month()))
}

func DateToUnix(date string) int64 {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", date, time.Local)
	return t.Unix()
}

func ISOMonth(sec, nsec int64) (int, int) {
	atime := UnixTime(sec, nsec)
	return atime.Year(), int(atime.Month())
}

func NewTimer(ms int) *time.Timer {
	return time.NewTimer(time.Millisecond * time.Duration(ms))
}

func NewTicker(ms int) *time.Ticker {
	return time.NewTicker(time.Millisecond * time.Duration(ms))
}

func After(ms int) <-chan time.Time {
	return time.After(time.Millisecond * time.Duration(ms))
}

func Tick(ms int) <-chan time.Time {
	return time.Tick(time.Millisecond * time.Duration(ms))
}

func Sleep(ms int) {
	time.Sleep(time.Millisecond * time.Duration(ms))
}

func SetTimeout(inteval int, fn func(...interface{}) int, args ...interface{}) {
	if inteval < 0 {
		LogError("new timerout inteval:%v", inteval)
		return
	}
	LogInfo("new timerout inteval:%v", inteval)

	Go2(func(cstop chan struct{}) {
		timeout := time.Millisecond * time.Duration(inteval)
		timer := time.NewTimer(timeout)
		for inteval > 0 {
			timeout = time.Millisecond * time.Duration(inteval)
			timer.Reset(timeout)
			select {
			case <-cstop:
				inteval = 0
			case <-timer.C:
				inteval = fn(args...)
			}
		}
		timer.Stop()
	})
}

func timerTick() {
	StartTick = time.Now().UnixNano() / 1000000
	NowTick = StartTick
	Timestamp = NowTick / 1000
	StartUnix = Timestamp
	Go(func() {
		for IsRuning() {
			Sleep(1)
			NowTick = time.Now().UnixNano() / 1000000
			Timestamp = NowTick / 1000
		}
	})
}

/**
* @brief 获得timestamp距离下个小时的时间，单位s
*
* @return uint32_t 距离下个小时的时间，单位s
 */
func GetNextHourIntervalS() int {
	return int(3600 - (Timestamp % 3600))
}

/**
 * @brief 获得timestamp距离下个小时的时间，单位ms
 *
 * @return uint32_t 距离下个小时的时间，单位ms
 */
func GetNextHourIntervalMS() int {
	return GetNextHourIntervalS() * 1000
}

/**
* @brief 时间戳转换为小时，24小时制，0点用24表示
*
* @param timestamp 时间戳
* @param timezone  时区
* @return uint32_t 小时 范围 1-24
 */
func GetHour24(timestamp int64, timezone int) int {
	hour := (int((timestamp%86400)/3600) + timezone)
	if hour > 24 {
		return hour - 24
	}
	return hour
}

/**
 * @brief 时间戳转换为小时，24小时制，0点用0表示
 *
 * @param timestamp 时间戳
 * @param timezone  时区
 * @return uint32_t 小时 范围 0-23
 */
func GetHour23(timestamp int64, timezone int) int {
	hour := GetHour24(timestamp, timezone)
	if hour == 24 {
		return 0 //24点就是0点
	}
	return hour
}

func GetHour(timestamp int64, timezone int) int {
	return GetHour23(timestamp, timezone)
}

/**
* @brief 判断两个时间戳是否是同一天
*
* @param now 需要比较的时间戳
* @param old 需要比较的时间戳
* @param timezone 时区
* @return uint32_t 返回不同的天数
 */
func IsDiffDay(now, old int64, timezone int) int {
	now += int64(timezone * 3600)
	old += int64(timezone * 3600)
	return int((now / 86400) - (old / 86400))
}

/**
* @brief 判断时间戳是否处于一个小时的两边，即一个时间错大于当前的hour，一个小于
*
* @param now 需要比较的时间戳
* @param old 需要比较的时间戳
* @param hour 小时，0-23
* @param timezone 时区
* @return bool true表示时间戳是否处于一个小时的两边
 */
func IsDiffHour(now, old int64, hour, timezone int) bool {
	diff := IsDiffDay(now, old, timezone)
	if diff == 1 {
		if GetHour23(old, timezone) >= hour {
			return GetHour23(now, timezone) >= hour
		} else {
			return true
		}
	} else if diff >= 2 {
		return true
	}

	return (GetHour23(now, timezone) >= hour) && (GetHour23(old, timezone) < hour)
}

/**
* @brief 判断时间戳是否处于跨周, 在周一跨天节点的两边
*
* @param now 需要比较的时间戳
* @param old 需要比较的时间戳
* @param hour 小时，0-23
* @param timezone 时区
* @return bool true表示时间戳是否处于跨周, 在周一跨天节点的两边
 */
func IsDiffWeek(now, old int64, hour, timezone int) bool {
	diffHour := IsDiffHour(now, old, hour, timezone)
	now += int64(timezone * 3600)
	old += int64(timezone * 3600)
	// 使用UTC才能在本地时间采用周一作为一周的开始
	_, nw := time.Unix(now, 0).UTC().ISOWeek()
	_, ow := time.Unix(old, 0).UTC().ISOWeek()
	return nw != ow && diffHour
}

// 当前任务周活跃6点重置
func IsDiffWeekSix(now, old int64, timezone int) bool {
	return SixUniqueWeek(now, 0) != SixUniqueWeek(old, 0)
}
func IsDiffMonthSix(now, old int64, timezone int) bool {
	return SixUniqueMonth(now, 0) != SixUniqueMonth(old, 0)
}

// 两个时间戳相差月份
func DiffMonth(t1, t2 int64) int {
	if t1 > t2 {
		t1, t2 = t2, t1
	}
	u1 := time.Unix(t1, 0)
	u2 := time.Unix(t2, 0)
	return (u2.Year()-u1.Year())*12 + int(u2.Month()) - int(u1.Month())
}

// 获取距离天数,相同的天返回1
func GetOffsetDay(time1 int64, time2 int64) int32 {
	t1 := ZeroTime(time1, 0) / SecondsPerDay
	t2 := ZeroTime(time2, 0) / SecondsPerDay
	return int32(math.Abs(float64(t1-t2))) + 1
}

// 计算指定时间戳到下个月一号0点的间隔秒数
func GetNextMonthInter(t1 int64) int64 {
	year := Year(t1, 0)
	month := Month(t1, 0)
	if month == 12 {
		year += 1
		month = 1
	} else {
		month += 1
	}
	nextMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	return nextMonth.Unix() - t1
}

// 计算指定时间戳到下个周星期一0点的间隔秒数
func GetNextWeekInter(t1 int64) int64 {
	day := 7 - WeekDayEx(t1, 0)                     //还剩多少整天
	inter := SecondsPerDay - (t1 - ZeroTime(t1, 0)) //今天还剩多少
	return int64(day)*SecondsPerDay + inter
}

// 获取下一个月一号0点的时间戳
func GetNextMonthUnix(sec int64) int64 {
	year, month, _ := time.Unix(sec, 0).Date()
	if month == 12 {
		year += 1
		month = 1
	} else {
		month += 1
	}
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).Unix()
}

func GetNextDayUnix(sec int64) int64 {
	year, month, day := time.Unix(sec, 0).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix() + 3600*24
}
