package utils

import (
	"log"
	"sync"
	"time"
)

func Timezone() (string, float32) {
	tz, offset := time.Now().Zone()
	return tz, float32(offset / 3600)
}

var locName2LocMap = map[string]*time.Location{}
var locName2LocMtx = sync.RWMutex{}

func LoadTimeLocation(name string) *time.Location {
	locName2LocMtx.RLock()
	if ret, ok := locName2LocMap[name]; ok {
		locName2LocMtx.RUnlock()
		return ret
	}
	locName2LocMtx.RUnlock()
	locName2LocMtx.Lock()
	if ret, ok := locName2LocMap[name]; ok {
		locName2LocMtx.Unlock()
		return ret
	}

	loc, err := time.LoadLocation(name)
	if err != nil {
		log.Printf("LoadTimeLocation %s err: %v", name, err)
		loc = time.UTC
	}
	locName2LocMap[name] = loc
	locName2LocMtx.Unlock()
	return loc
}

func IsSameDay(t1, t2 int64) bool {
	tt1 := time.Unix(t1, 0)
	tt2 := time.Unix(t2, 0)
	y1, m1, d1 := tt1.Date()
	y2, m2, d2 := tt2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func IsSameDayWithLocation(t1, t2 int64, loc *time.Location) bool {
	tt1 := time.Unix(t1, 0).In(loc)
	tt2 := time.Unix(t2, 0).In(loc)
	y1, m1, d1 := tt1.Date()
	y2, m2, d2 := tt2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func IsSameWeek(t1, t2 int64) bool {
	tt1 := time.Unix(t1, 0)
	tt2 := time.Unix(t2, 0)
	y1, w1 := tt1.ISOWeek()
	y2, w2 := tt2.ISOWeek()
	return y1 == y2 && w1 == w2
}

func IsSameWeekWithLocation(t1, t2 int64, loc *time.Location) bool {
	tt1 := time.Unix(t1, 0).In(loc)
	tt2 := time.Unix(t2, 0).In(loc)
	y1, w1 := tt1.ISOWeek()
	y2, w2 := tt2.ISOWeek()
	return y1 == y2 && w1 == w2
}

func IsSameMonth(t1, t2 int64) bool {
	tt1 := time.Unix(t1, 0)
	tt2 := time.Unix(t2, 0)
	y1, m1, _ := tt1.Date()
	y2, m2, _ := tt2.Date()
	return y1 == y2 && m1 == m2
}

func IsSameMonthWithLocation(t1, t2 int64, loc *time.Location) bool {
	tt1 := time.Unix(t1, 0).In(loc)
	tt2 := time.Unix(t2, 0).In(loc)

	y1, m1, _ := tt1.Date()
	y2, m2, _ := tt2.Date()
	return y1 == y2 && m1 == m2
}
