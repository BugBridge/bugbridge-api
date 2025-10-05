package util

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func BuildUpdate(v any) bson.M {
	update := bson.M{}
	buildUpdate("", reflect.ValueOf(v), update)
	return update
}

func buildUpdate(prefix string, v reflect.Value, update bson.M) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		key := strings.Split(tag, ",")[0]
		if key == "" {
			continue
		}

		// full key = parent.child if prefix exists
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if value.Kind() == reflect.Ptr && value.IsNil() {
			continue
		}

		// handle nested structs (not time.Time)
		if value.Kind() == reflect.Struct && value.Type().String() != "time.Time" {
			buildUpdate(fullKey, value, update)
			continue
		}

		// skip zero values
		if isZero(value) {
			continue
		}

		update[fullKey] = value.Interface()
	}
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.IsNil()
	case reflect.Struct:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	default:
		return false
	}
}
