package api

import (
	"strconv"
	"time"
)

type DateTime time.Time

func (dt *DateTime) MarshalJSON() ([]byte, error) {
	ts := time.Time(*dt).Unix()
	return []byte(strconv.FormatInt(ts, 10)), nil
}

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*dt = DateTime(time.Unix(int64(ts), 0))
	return nil
}
