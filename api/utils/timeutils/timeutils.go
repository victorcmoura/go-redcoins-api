package timeutils

import (
	"time"
)

func BeginningOfTheDay(t time.Time) time.Time {
    year, month, day := t.Date()
    return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndOfTheDay(t time.Time) time.Time {
    return BeginningOfTheDay(t).Add(time.Hour * 24)
}