package util

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// UrlFormEncodedToObject converts a url form encoded string to a struct
func UrlFormEncodedToObject(str string, obj any) error {
	values, err := url.ParseQuery(str)
	if err != nil {
		return err
	}

	objValue := reflect.ValueOf(obj).Elem()
	objType := objValue.Type()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := objValue.FieldByName(field.Name)
		fieldTag := field.Tag.Get("json")

		if strings.Contains(fieldTag, ",") {
			fieldTag = strings.Split(fieldTag, ",")[0]
		}

		if values.Get(fieldTag) != "" && fieldValue.IsValid() && fieldValue.CanSet() {
			fieldType := fieldValue.Type()
			switch fieldType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v, err := strconv.ParseInt(values.Get(fieldTag), 10, 64)
				if err != nil {
					return err
				}
				fieldValue.SetInt(v)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				v, err := strconv.ParseUint(values.Get(fieldTag), 10, 64)
				if err != nil {
					return err
				}
				fieldValue.SetUint(v)
			case reflect.Float32, reflect.Float64:
				v, err := strconv.ParseFloat(values.Get(fieldTag), 64)
				if err != nil {
					return err
				}
				fieldValue.SetFloat(v)
			case reflect.Bool:
				v, err := strconv.ParseBool(values.Get(fieldTag))
				if err != nil {
					return err
				}
				fieldValue.SetBool(v)
			case reflect.String:
				fieldValue.SetString(values.Get(fieldTag))
			}
		}
	}

	return nil
}

func ObjectToUrlFormEncoded(obj any) ([]byte, error) {
	objValue := reflect.ValueOf(obj).Elem()
	objType := objValue.Type()
	values := url.Values{}
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := objValue.FieldByName(field.Name)
		fieldTag := field.Tag.Get("json")

		if strings.Contains(fieldTag, ",") {
			split := strings.Split(fieldTag, ",")

			if fieldValue.Kind() == reflect.String && fieldValue.String() == "" {
				if strings.Contains(split[1], "omitempty") {
					continue
				}
			}
			fieldTag = split[0]
		}

		if fieldValue.IsValid() && fieldValue.CanSet() {
			fieldType := fieldValue.Type()
			switch fieldType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				values.Set(fieldTag, strconv.FormatInt(fieldValue.Int(), 10))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				values.Set(fieldTag, strconv.FormatUint(fieldValue.Uint(), 10))
			case reflect.Float32, reflect.Float64:
				values.Set(fieldTag, strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64))
			case reflect.Bool:
				values.Set(fieldTag, strconv.FormatBool(fieldValue.Bool()))
			case reflect.String:
				values.Set(fieldTag, fieldValue.String())
			}
		}
	}

	return []byte(values.Encode()), nil
}
