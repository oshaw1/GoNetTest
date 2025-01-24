package dataManagement

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
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

// Gets all test date directories sorted newest first
func (r *Repository) GetTestDirectories() ([]string, error) {
	files, err := os.ReadDir(r.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var dates []string
	for _, f := range files {
		if f.IsDir() {
			dates = append(dates, f.Name())
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(dates)))
	return dates, nil
}

func (r *Repository) GetTestFilesInGroup(date, testType string) (map[string][]string, error) {
	testPath := filepath.Join(r.baseDir, date, testType)
	files, err := os.ReadDir(testPath)
	if err != nil {
		return nil, err
	}

	groups := make(map[string][]string)
	for _, file := range files {
		name := file.Name()
		timeGroup := strings.Split(name, "_")[2]
		groups[timeGroup] = append(groups[timeGroup], filepath.Join(testPath, name))
	}
	return groups, nil
}

func (r *Repository) ListTestTypesInDateDir(date string) ([]string, error) {
	dateDir := filepath.Join(r.baseDir, date)
	files, err := os.ReadDir(dateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var types []string
	for _, f := range files {
		if f.IsDir() {
			types = append(types, f.Name())
		}
	}

	sort.Strings(types)
	return types, nil
}
