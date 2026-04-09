package binding

import (
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

type defaultValidator struct {
	validate *validator.Validate
	once     sync.Once
}

var _ StructValidator = &defaultValidator{}

func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
	})
}

func (v *defaultValidator) Engine() any {
	v.lazyInit()
	return v.validate
}

func (v *defaultValidator) ValidateStruct(obj any) error {
	// Алгоритм:
	// 1. Узнать тип объекта через reflect.ValueOf(obj)
	val := reflect.ValueOf(obj)

	// 2. Если это указатель — получить значение, на которое он указывает (.Elem())
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// 3. Если это структура (.Kind() == reflect.Struct) — вызвать v.lazyInit() и v.validate.Struct(obj)
	if val.Kind() == reflect.Struct {
		v.lazyInit()
		err := v.validate.Struct(obj)
		if err != nil {
			return err
		}
	}

	// 4. Если это не структура — вернуть nil
	return nil
}
