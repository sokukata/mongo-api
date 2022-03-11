package convert

import (
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertDoc(d primitive.D) (map[string]interface{}, error) {
	m := make(map[string]interface{}, len(d))
	for _, e := range d {
		value, err := convertElem(e.Value)
		if err != nil {
			return nil, err
		}
		// Add to new map
		m[e.Key] = value
	}
	return m, nil
}

// Parse one field type
func convertElem(val interface{}) (interface{}, error) {
	// nil is JSON null. just ignore this
	if val == nil {
		return nil, nil
	}

	switch val := val.(type) {

	case string, bool, int32, int64:
		// Keep native type
		return val, nil

	case float64:
		// Avoid values invalid for json
		if math.IsInf(val, 0) || math.IsNaN(val) {
			val = 0
		}
		return val, nil

	case primitive.DateTime:
		// Bson datetime is number of millisec since epoch
		dt := int64(val)
		if dt == 0 {
			return time.Time{}, nil
		}
		tm := time.Unix(dt/1000, dt%1000*1000000).UTC()
		return tm, nil

	case primitive.ObjectID:
		// ObjectID ("internal primary key"): just convert to string
		str := val.Hex()
		return str, nil

	case primitive.D:
		return ConvertDoc(val)
	case primitive.A:
		// convert array values
		newArray := make([]interface{}, len(val))
		for i, v := range val {
			_, isDoc := v.(primitive.D)
			v, err := convertElem(v)
			if err != nil {
				return nil, err
			}
			if isDoc && v == nil {
				// Avoid nil inside arrays if elements are supposed to be documents
				v = make(map[string]interface{})
			}
			newArray[i] = v
		}
		return newArray, nil
	}

	return nil, fmt.Errorf("Cannot convert value from MongoDB (unsupported data type '%T' for value '%v')", val, val)
}
