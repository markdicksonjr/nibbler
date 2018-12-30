package nibbler

import (
	"encoding/json"
	"reflect"
)

// typeOfPtr should be the type of a pointer to the type you're unmarshalling to
func FromJson(jsonString string, typeOfPtr reflect.Type) (interface{}, error) {
	byteArray := []byte(jsonString)

	objType := typeOfPtr.Elem()
	obj := reflect.New(objType).Interface()

	err := json.Unmarshal(byteArray, obj)
	return obj, err
}

func ToJson(obj interface{}) (result string, err error) {
	objJsonBytes, err := json.Marshal(obj)

	if err != nil {
		return
	}

	result = string(objJsonBytes)
	return
}
