package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const layout = "2006-01-02 15:04:05.999999999-07:00"

type Timestamp time.Time

func (st *Timestamp) Scan(value any) error {
	if value == nil {
		*st = Timestamp(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		// go-sqlite3 already parses columns declared as TIMESTAMP/DATETIME/DATE.
		*st = Timestamp(v)
		return nil
	case string:
		t, err := time.Parse(layout, v)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %s", v)
		}
		*st = Timestamp(t)
		return nil
	default:
		return fmt.Errorf("expected string or time.Time, got %T", value)
	}
}

func (st Timestamp) Value() (driver.Value, error) {
	t := time.Time(st)

	if t.IsZero() {
		return nil, nil
	}

	return t.Format(layout), nil
}

func (st Timestamp) Time() time.Time {
	return time.Time(st)
}
