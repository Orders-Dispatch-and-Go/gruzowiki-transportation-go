package duration

import (
	"time"

	"github.com/goccy/go-json"
)

type Seconds struct {
	time.Duration
}

func (d *Seconds) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Seconds())
}

func (d *Seconds) UnmarshalJSON(b []byte) error {
	var seconds int64
	if err := json.Unmarshal(b, &seconds); err != nil {
		return err
	}

	d.Duration = time.Duration(seconds * int64(time.Second))

	return nil
}
