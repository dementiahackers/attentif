package entry

import (
	"database/sql"
	"time"
)

type Entry struct {
	ID          int64
	UserID      int64
	Rate        int64
	Description sql.NullString
	CreatedAt   time.Time
}

func (e *Entry) Timestamp() string {
	return e.CreatedAt.Local().Format("Jan 2, 3:04pm")
}

func (e *Entry) Feeling() string {
	if e.Rate < 33 {
		return "sad"
	} else if e.Rate < 66 {
		return "neutral"
	} else {
		return "happy"
	}
}

func GroupByRating(entries []Entry) []int {
	r := make([]int, 11)
	for _, e := range entries {
		i := e.Rate % 10
		if i == 0 && e.Rate == 100 {
			i = 10
		}
		r[i]++
	}
	return r
}
