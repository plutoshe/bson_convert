package bson_convert

import (
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

type bsonPattern struct {
	required, ignored map[string]struct{}
	omitzero          bool
}

// required first, ignored later.
func Convert2BSON(source interface{}, omitzero bool, requiredField []string, ignoredField []string) bson.M {
	required := make(map[string]struct{})
	ignored := make(map[string]struct{})
	for _, v := range requiredField {
		required[v] = struct{}{}
	}
	for _, v := range ignoredField {
		if _, require := required[v]; !require {
			ignored[v] = struct{}{}
		}
	}
	b := bsonPattern{required, ignored, omitzero}
	return b.ConvertToBSON(source)
}

func (b *bsonPattern) ConvertToBSON(source interface{}) bson.M {
	sourceType := reflect.TypeOf(source)
	v := reflect.ValueOf(source)
	res := bson.M{}
	for i := 0; i < v.NumField(); i++ {
		if sourceType.Field(i).Name[0] >= 'A' && sourceType.Field(i).Name[0] <= 'Z' {
			tag := sourceType.Field(i).Tag.Get("bson")
			fields := strings.Split(tag, ",")
			if len(fields) > 1 {
				tag = fields[0]
			}
			if _, ignore := b.ignored[tag]; ignore {
				continue
			}
			if _, require := b.required[tag]; !require && isEmptyValue(v.Field(i)) && b.omitzero {
				continue
			}
			res[tag] = b.getValue(v.Field(i))
		}
	}
	return res
}

func (b *bsonPattern) getValue(v reflect.Value) interface{} {
	if v.Kind() == reflect.Map || v.Kind() == reflect.Struct {
		return b.ConvertToBSON(v.Interface())
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
