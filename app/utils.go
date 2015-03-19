package app

import (
	"fmt"
	"time"

	"github.com/revel/revel"
)

// assert
func AssertNoError(err error, msg string) {
	if err != nil {
		Assert(false, msg+"("+err.Error()+")")
	}
}

func Assert(assertion bool, msg string) {
	if !assertion {
		err := "assertion failed"
		if msg != "" {
			err += ": " + msg
		}
		panic(err)
	}
}

// time.layout
const (
	DefaultDateTimeFull = "2016-01-02 15:04:05.999"
	DefaultDateTime     = "2016-01-02 15:04:05"
	DefaultDate         = "2016-01-02"
	DefaultTime         = "15:04:05"
)

// time.format
// FormatDefault returns time as DefaultDateTime format
func FormatDefault(t time.Time) string {
	return time.Time.Format(t, DefaultDateTime)
}

type RichTime time.Time

// Yesterday for short
func (r RichTime) Yesterday() time.Time {
	var t = time.Time(r)
	return t.AddDate(0, 0, -1)
}

// revel
// Get config from revel, panic if not exists
func ForceGetConfig(key string) string {
	v, exists := revel.Config.String(key)
	Assert(exists, fmt.Sprintf("Missing revel app config for key:%s", key))
	// fmt.Printf("[%s]: %v", key, v)
	return v
}
