package convert

import (
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const KEY_NAME = "_id"

func ConvertMap(m map[string]interface{}, update bool) bson.M {
	cp := make(bson.M, len(m))
	for k, v := range m {
		if update && k == KEY_NAME {
			continue
		}
		val := convert(v, update)
		cp[k] = val
	}
	if update {
		// Wrap packet
		c := make(bson.M)
		c["$set"] = cp
		return c
	}
	return cp
}

func convertSlice(a []interface{}) bson.A {
	cp := make(bson.A, len(a))
	for i, v := range a {
		cp[i] = convert(v, false)
	}
	return cp
}

func convert(v interface{}, update bool) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		return ConvertMap(v, update)
	case []interface{}:
		return convertSlice(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
		if dt, err := time.Parse("2006-01-02T15:04:05 -07:00", v); err == nil {
			return dt
		}
	}
	// bool, string:
	return v
}
