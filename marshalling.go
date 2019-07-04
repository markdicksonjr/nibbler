package nibbler

import (
	"encoding/json"
	"reflect"
)

// typeOfPtr should be the type of a pointer to the type you're unmarshalling to
func FromJson(jsonString string, typeItem reflect.Type) (interface{}, error) {
	obj := reflect.New(typeItem).Interface()
	err := json.Unmarshal([]byte(jsonString), obj)
	return obj, err
}

func ToJson(obj interface{}) (result []byte, err error) {
	return json.Marshal(obj)
}
