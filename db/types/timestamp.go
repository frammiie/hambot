package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Timestamp time.Time

func (st *Timestamp) Scan(value any) error {
	if value == nil {
		*st = Timestamp(time.Time{})
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}

	t, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", str)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %s", str)
	}

	*st = Timestamp(t)

	return nil
}

func (st Timestamp) Value() (driver.Value, error) {
	t := time.Time(st)

	if t.IsZero() {
		return nil, nil
	}
	
	return t.Format("2006-01-02 15:04:05.999999999-07:00"), nil
}

func (st Timestamp) Time() time.Time {
	return time.Time(st)
}
