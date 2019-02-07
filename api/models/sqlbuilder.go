package models

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)


func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	case reflect.TypeOf(time.Time{}).Kind():
		return reflect.ValueOf(time.Time{}).Interface() == v.Interface()
	}
	return false
}

func createSet(obj interface{}) string {
	rt := reflect.TypeOf(obj) // reflect.Type
	rv := reflect.ValueOf(obj) // reflect.Value

	setQuery := "SET"
	for i := 0; i < rv.NumField(); i++ {
		if !isZero(rv.Field(i)) {
			tagValue := strings.SplitN(rt.Field(i).Tag.Get("json"), ",", 2)[0]
			value := rv.Field(i).Interface()
			if rv.Field(i).Kind() == reflect.String {
				value = fmt.Sprintf("'%s'", value)
			}
			if i > 0 {
				setQuery += " AND"
			}
			setQuery += fmt.Sprintf(" %s = %s", tagValue, value)

		}
	}
	return setQuery
}