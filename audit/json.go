package audit

import "encoding/json"

// TryMarshall attempts to marshal the given object into a JSON string.
// If the object cannot be marshalled, nil is returned.
func TryMarshall(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return data
}
