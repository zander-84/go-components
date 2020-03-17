package CHelperTime

import (
	"sync"
	"time"
)

var _timeZone = new(TimeZone)

type TimeZone struct {
	location *time.Location
	mutex    sync.Mutex
	once     sync.Once
}

func NewTimeZone(timezZone string) interface{} {
	_timeZone.once.Do(func() {
		if timezZone == "" {
			timezZone = "Asia/Shanghai"
		}
		_timeZone.location, _ = time.LoadLocation(timezZone)
	})

	return _timeZone
}

func (t *TimeZone) SetLocation(loc string) (err error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	_loc, err := time.LoadLocation(loc)
	if err == nil {
		t.location = _loc
	}
	return err
}

func (t *TimeZone) Location() *time.Location {
	return t.location
}

func (t *TimeZone) Now() time.Time {
	return time.Now().In(t.location)
}

func (t *TimeZone) LocationNow(src time.Time) time.Time {
	return src.In(t.location)
}

func (t *TimeZone) Format(layout string) string {
	return t.Now().Format(layout) //"2006-01-02 15:04:05"
}

// "01/02/2006", "02/08/2015"
func (t *TimeZone) Parse(layout string, sourceTime string) (time.Time, error) {
	t1, err := time.ParseInLocation(layout, sourceTime, t.location)
	return t1.In(t.location), err
}

func (t *TimeZone) Year() string {
	return t.Now().Format("2006")
}

func (t *TimeZone) Month() string {
	return t.Now().Format("01")
}

func (t *TimeZone) Day() string {
	return t.Now().Format("02")
}

func (t *TimeZone) FormatSlash() string {
	return t.Now().Format("2006/01/02 15:04:05")
}

func (t *TimeZone) FormatHyphen() string {
	return t.Now().Format("2006-01-02 15:04:05")
}
