package util

import (
	"time"
)

func ToTimestamp(s string) int64 {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}

	ts := t.Unix()
	return ts
}
