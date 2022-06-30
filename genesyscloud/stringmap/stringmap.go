package stringmap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetNillableValue[T any](m map[string]interface{}, key string) *T {
	value, ok := m[key]
	if ok {
		v := value.(T)
		return &v
	}
	return nil
}

func GetNonDefaultValue[T comparable](m map[string]interface{}, key string) *T {
	value := GetNillableValue[T](m, key)
	if value != nil {
		defaultValue := new(T)
		if *value != *defaultValue {
			return value
		}
	}
	return nil
}

func SetValueIfNotNil[T any](m map[string]interface{}, key string, value *T) {
	if value != nil {
		m[key] = *value
	}
}

func BuildSdkStringList(d map[string]interface{}, attrName string) *[]string {
	return BuildSdkList[string](d, attrName, nil)
}

func BuildSdkList[T interface{}](d map[string]interface{}, attrName string, elementBuilder func(map[string]interface{}) *T) *[]T {
	child := d[attrName]
	if child != nil {
		list := child.(*schema.Set).List()
		sdkList := make([]T, len(list))
		for i, element := range list {
			switch element.(type) {
			case T:
				sdkList[i] = element.(T)
			case map[string]interface{}:
				sdkList[i] = *elementBuilder(element.(map[string]interface{}))
			}
		}
		return &sdkList
	}
	return nil
}