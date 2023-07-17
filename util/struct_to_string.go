package util

import (
  "reflect"
  "fmt"
)

func StructToStringFormat(obj any) string {
	result := ""
	v := reflect.ValueOf(obj)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		result += fmt.Sprintf("%s=%v;", field.Name, value)
	}

	// Remove the trailing semicolon if present
	if len(result) > 0 && result[len(result)-1] == ';' {
		result = result[:len(result)-1]
	}

	return result
}
