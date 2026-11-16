package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

func ApplyDefaultValues(struct_ any) error {
	v := reflect.ValueOf(struct_).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field, structField := v.Field(i), t.Field(i)
		if !field.CanSet() {
			continue
		}
		defaultValue := structField.Tag.Get("default")
		if defaultValue == "" {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(defaultValue)
		case reflect.Bool:
			switch defaultValue {
			case "true":
				field.SetBool(true)
			case "false":
				field.SetBool(false)
			default:
				return fmt.Errorf("invalid bool expression: %v, use true/false", defaultValue)
			}
		case reflect.Int:
			val, err := strconv.Atoi(defaultValue)
			if err != nil {
				return err
			}
			field.SetInt(int64(val))
		}
	}
	return nil
}
