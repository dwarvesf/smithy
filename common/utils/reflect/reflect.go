package reflect

import (
	"errors"
	"reflect"
)

func ConvertFromInterfacePtr(v interface{}) (interface{}, error) {
	var value interface{}
	it := reflect.ValueOf(v).Elem().Interface()
	typee := reflect.ValueOf(it).Kind().String()
	switch typee {
	case "int64":
		value = int(reflect.ValueOf(v).Elem().Interface().(int64))
	case "float64":
		value = reflect.ValueOf(v).Elem().Interface().(float64)
	case "string":
		value = reflect.ValueOf(v).Elem().Interface().(string)
	default:
		return nil, errors.New("type is not supported")
	}
	return value, nil
}
