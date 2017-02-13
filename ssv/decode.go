package ssv

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/0xdeafcafe/go-xbdm/helpers"
)

// Unmarshal parses the SSV-encoded data and stores the result
// in the value pointed to by v.
func Unmarshal(str string, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("v must be a pointer and not nil")
	}

	rt := reflect.Indirect(rv).Elem().Type()
	rv = reflect.Indirect(rv).Elem()
	ssvData := helpers.ParseSpaceSeparatedValues(str)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.FieldByName(field.Name)
		fieldKind := field.Type.Kind()

		xbdmTag, ok := field.Tag.Lookup("goxbdm")
		if !ok {
			errStr := fmt.Sprintf("field %s is missing a required tag", field.Name)
			return errors.New(errStr)
		}

		value, ok := ssvData[xbdmTag]
		if !ok {
			value = ""
		}

		if !fieldValue.CanSet() {
			continue
		}

		switch fieldKind {
		case reflect.String:
			fieldValue.SetString(value)
			break
		case reflect.Int64:
			val := helpers.ConvertHexToInt64(value)
			fieldValue.SetInt(val)
			break
		default:
			errStr := fmt.Sprintf("the kind %s cannot be unmarshalled", fieldKind.String())
			return errors.New(errStr)
		}
	}

	return nil
}
