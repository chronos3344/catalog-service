package util

import (
	"time"
)

// Duration - кастомный тип для парсинга времени из строки
type Duration struct {
	time.Duration
}

// UnmarshalText реализует интерфейс encoding.TextUnmarshaler
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
