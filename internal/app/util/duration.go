package util

import (
	//"fmt"
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

// MarshalText реализует интерфейс encoding.TextMarshaler
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// String возвращает строковое представление
func (d Duration) String() string {
	return d.Duration.String()
}
