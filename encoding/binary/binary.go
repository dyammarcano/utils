package binary

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

var endian = binary.LittleEndian

func Marshal(data any) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := serializeStruct(buf, reflect.ValueOf(data)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func serializeStruct(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if err := serializeStruct(buf, v.Field(i)); err != nil {
				return err
			}
		}
	case reflect.Int32, reflect.Float64, reflect.Bool:
		if err := binary.Write(buf, endian, v.Interface()); err != nil {
			return err
		}
	case reflect.String:
		str := v.String()
		if err := binary.Write(buf, endian, int32(len(str))); err != nil {
			return err
		}
		if err := binary.Write(buf, endian, []byte(str)); err != nil {
			return err
		}
	}
	return nil
}

func Unmarshal(data []byte, ptr any) error {
	buf := bytes.NewReader(data)
	return deserializeStruct(buf, reflect.ValueOf(ptr).Elem())
}

func deserializeStruct(buf *bytes.Reader, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if err := deserializeStruct(buf, field); err != nil {
				return err
			}
		}
	case reflect.Int32, reflect.Float64, reflect.Bool:
		if err := binary.Read(buf, endian, v.Addr().Interface()); err != nil {
			return err
		}
	case reflect.String:
		var strLen int32
		if err := binary.Read(buf, endian, &strLen); err != nil {
			return err
		}
		strBytes := make([]byte, strLen)
		if err := binary.Read(buf, endian, &strBytes); err != nil {
			return err
		}
		v.SetString(string(strBytes))
	}
	return nil
}
