package dataManagement

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"
)

func (r *Repository) readJSONFile(path string) (interface{}, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var result interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return result, nil
}

func (r *Repository) convertToMap(data interface{}) (map[string]interface{}, error) {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() == reflect.Struct {
		return structToMap(value), nil
	}

	if value.Kind() == reflect.Map {
		if m, ok := data.(map[string]interface{}); ok {
			return m, nil
		}
	}

	return nil, fmt.Errorf("unsupported data type: %v", value.Kind())
}

func structToMap(value reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		fieldValue := value.Field(i)
		result[field.Name] = fieldValue.Interface()
	}
	return result
}

func isDateFolder(name string) bool {
	_, err := time.Parse(dateFormat, name)
	return err == nil
}
