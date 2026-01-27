package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	"github.com/chronos3344/catalog-service/internal/app/util"
	"github.com/joho/godotenv"
)

var Root = struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}{}

func Load() {
	_ = godotenv.Load(".env")
	if err := load(&Root, "APP"); err != nil {
		log.Fatalf("Config error: %v", err)
	}
	log.Println("Config loaded")
}

func load(config interface{}, prefix string) error {
	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)

		if f.Kind() == reflect.Struct {
			p := ft.Tag.Get("env")
			if p == "" {
				p = prefix + "_" + strings.ToUpper(ft.Name)
			}
			if err := load(f.Addr().Interface(), p); err != nil {
				return err
			}
			continue
		}

		key := ft.Tag.Get("env")
		if key == "" {
			key = prefix + "_" + strings.ToUpper(ft.Name)
		}

		val := os.Getenv(key)
		if val == "" {
			if strings.Contains(ft.Tag.Get("validate"), "required") {
				return missingField(key)
			}
			continue
		}

		if err := set(&f, val); err != nil {
			return fieldError(ft.Name, err)
		}
	}
	return nil
}

func set(f *reflect.Value, val string) error {
	switch f.Kind() {
	case reflect.String:
		f.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setInt(f, val)
	default:
		if f.CanInterface() {
			if u, ok := f.Interface().(interface{ UnmarshalText(data []byte) error }); ok {
				return u.UnmarshalText([]byte(val))
			}
		}
		return unsupportedType(f.Type().String())
	}
	return nil
}

func setInt(f *reflect.Value, val string) error {
	if f.Type().String() == "util.Duration" {
		var d util.Duration
		if err := d.UnmarshalText([]byte(val)); err != nil {
			return err
		}
		f.Set(reflect.ValueOf(d))
		return nil
	}

	n, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return parseError(val)
	}
	f.SetInt(n)
	return nil
}

// Минимальные вспомогательные функции
func missingField(key string) error {
	return errorf("missing required: %s", key)
}

func fieldError(name string, err error) error {
	return errorf("field %s: %w", name, err)
}

func unsupportedType(typ string) error {
	return errorf("unsupported type: %s", typ)
}

func parseError(val string) error {
	return errorf("parse error: %s", val)
}

func errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
