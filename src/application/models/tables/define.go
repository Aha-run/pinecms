package tables

import (
	"database/sql/driver"
	"time"
)

const localDateTimeFormat string = "2006-01-02 15:04:05"

type LocalTime time.Time

func (l LocalTime) MarshalJSON() ([]byte, error) {
	if time.Time(l).IsZero() {
		return []byte(`""`), nil
	}
	b := make([]byte, 0, len(localDateTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(l).AppendFormat(b, localDateTimeFormat)
	b = append(b, '"')
	return b, nil
}

func (l *LocalTime) UnmarshalJSON(b []byte) error {
	now, err := time.ParseInLocation(`"`+localDateTimeFormat+`"`, string(b), time.Local)
	*l = LocalTime(now)
	return err
}

func (l LocalTime) String() string {
	return time.Time(l).Format(localDateTimeFormat)
}

func (l LocalTime) Now() LocalTime {
	return LocalTime(time.Now())
}

func (l LocalTime) ParseTime(t time.Time) LocalTime {
	return LocalTime(t)
}

func (l LocalTime) format() string {
	return time.Time(l).Format(localDateTimeFormat)
}

func (l LocalTime) MarshalText() ([]byte, error) {
	return []byte(l.format()), nil
}

func (l *LocalTime) FromDB(b []byte) error {
	if nil == b || len(b) == 0 {
		return nil
	}
	var now time.Time
	var err error
	if now, err = time.ParseInLocation(localDateTimeFormat, string(b), time.Local); err != nil {
		return err
	}
	*l = LocalTime(now)
	return nil
}

func (l *LocalTime) ToDB() ([]byte, error) {
	if nil == l {
		return nil, nil
	}
	return []byte(time.Time(*l).Format(localDateTimeFormat)), nil
}

func (l *LocalTime) Value() (driver.Value, error) {
	if nil == l {
		return nil, nil
	}
	return time.Time(*l).Format(localDateTimeFormat), nil
}
