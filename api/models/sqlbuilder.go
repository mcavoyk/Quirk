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
	case reflect.Struct:
		structVal := v.Interface()
		switch t := structVal.(type) {
		case Default:
			return true
		case time.Time:
			return time.Time{} == t
		}
	}

	return false
}

func UpdateValues(table string, obj interface{}) string {
	return fmt.Sprintf("UPDATE %s %s WHERE id = ?", table, createSet(obj))
}

func createSet(obj interface{}) string {
	rt := reflect.TypeOf(obj)  // reflect.Type
	rv := reflect.ValueOf(obj) // reflect.Value

	setQuery := "SET"
	fieldIndex := 0
	for i := 0; i < rv.NumField(); i++ {
		if !isZero(rv.Field(i)) {
			tagValue := strings.SplitN(rt.Field(i).Tag.Get("json"), ",", 2)[0]
			value := rv.Field(i).Interface()
			if rv.Field(i).Kind() == reflect.String {
				value = fmt.Sprintf("'%s'", value)
			}
			if fieldIndex > 0 {
				setQuery += ", "
			}

			setQuery += fmt.Sprintf(" %s = %v", tagValue, value)
			fieldIndex++
		}
	}
	return setQuery
}

func InsertValues(insert string) string {
	start := strings.Split(insert, "(")
	if len(start) == 1 {
		return insert
	}
	start[1] = strings.Join(start[1:], "(")

	end := strings.Split(start[1], ")")
	endStmt := ""
	if len(end) > 1 {
		endStmt = strings.Join(end[1:], ")")
	}

	columns := end[0]
	columnSplit := strings.Split(columns, ",")
	namedValues := ""
	for i, column := range columnSplit {
		column = strings.TrimSpace(column)
		if i != 0 {
			namedValues += ", "
		}
		namedValues += fmt.Sprintf(":%s", column)
	}

	return fmt.Sprintf("%s(%s) VALUES (%s)%s", start[0], end[0], namedValues, endStmt)
}
