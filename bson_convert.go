package bson_convert

import (
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// ignore first, required later.
func Convert2BSON(source interface{}, omitzero bool, requiredField []string, ignoredField []string) bson.M {
	v := reflect.ValueOf(source)
	sourceType := reflect.TypeOf(source)

	required := make(map[string]struct{})
	ignored := make(map[string]struct{})
	for _, v := range requiredField {
		required[v] = struct{}{}
	}
	for _, v := range ignoredField {
		ignored[v] = struct{}{}
	}
	res := bson.M{}

	for i := 0; i < v.NumField(); i++ {

		tag := sourceType.Field(i).Tag.Get("bson")
		fields := strings.Split(tag, ",")
		if len(fields) > 1 {
			tag = fields[0]
		}
		if _, ignore := ignored[tag]; ignore {
			continue
		}
		if _, require := required[tag]; !require && isEmptyValue(v.Field(i)) && omitzero {
			continue
		}
		res[tag] = getValue(v.Field(i), omitzero, requiredField, ignoredField)
	}
	return res
}

func getValue(v reflect.Value, omitzero bool, requiredField []string, ignoredField []string) interface{} {
	if v.Kind() == reflect.Map || v.Kind() == reflect.Struct {
		return Convert2BSON(v.Interface(), omitzero, requiredField, ignoredField)
	}
	return v.Interface()
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
