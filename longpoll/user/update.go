package longpoll

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

// Update represents the information an event contains.
type Update struct {
	Type    int64
	Data    interface{}
	RawData json.RawMessage
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (update *Update) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &update.Data); err != nil {
		return err
	}

	update.Type = int64(update.Data.([]interface{})[0].(float64))
	update.RawData = data

	switch update.Type {
	case EventNewMessage:
		update.Data = &NewMessage{}
	default:
		return nil // []interface{} remains the type of update.Data
	}

	return update.Unmarshal(update.Data)
}

// Unmarshal parses JSON-encoded update.RawData and stores the result
// in the struct value pointed to by outputStruct.
func (update *Update) Unmarshal(outputStruct interface{}) error {
	structPtr := reflect.ValueOf(outputStruct)

	if structPtr.Kind() != reflect.Ptr || structPtr.IsNil() || structPtr.Elem().Kind() != reflect.Struct {
		return errors.New("outputStruct must be a valid pointer to a struct")
	}

	structValue := structPtr.Elem()

	arr := make([]interface{}, structValue.NumField()+1)

	for n := 0; n < structValue.NumField(); n++ {
		arr[n+1] = structValue.Field(n).Addr().Interface() // skip event code
	}

	return json.Unmarshal(update.RawData, &arr)
}
