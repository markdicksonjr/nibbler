package nibbler

import (
	"encoding/json"
	"reflect"
)

// typeOfPtr should be the type of a pointer to the type you're unmarshalling to
func FromJson(jsonString string, typeItem reflect.Type) (interface{}, error) {
	byteArray := []byte(jsonString)

	obj := reflect.New(typeItem).Interface()

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
