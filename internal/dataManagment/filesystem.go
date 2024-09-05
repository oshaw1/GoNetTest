package dataManagment

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/oshaw1/go-net-test/config"
)

const dateFormat = "2006-01-02"

func SaveTestData(data interface{}, test string) error {
	now := time.Now()

	dir, err := generateFilePath("data/output/", now, test)
	if err != nil {
		return err
	}

	filename := getUniqueFilename(test, dir, now)
	fullPath := filepath.Join(dir, filename)

	// Check if data is a pointer and get its underlying value
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Convert the data to a map
	var dataMap map[string]interface{}
	if value.Kind() == reflect.Struct {
		dataMap = structToMap(value)
	} else if value.Kind() == reflect.Map {
		dataMap = data.(map[string]interface{})
	} else {
		return fmt.Errorf("unsupported data type: %v", value.Kind())
	}

	jsonData, err := json.MarshalIndent(dataMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	err = os.WriteFile(fullPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	fmt.Printf("Data saved successfully to: %s\n", fullPath)
	return nil
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

func CheckForRecentTestData(rootDir string, fileExtension string) (bool, string, error) {
	conf, err := config.NewConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	cutoffTime := time.Now().AddDate(0, 0, -conf.RecentDays)
	return walkRootDirectory(rootDir, cutoffTime, fileExtension)
}

func walkRootDirectory(rootDir string, cutoffTime time.Time, fileExtension string) (bool, string, error) {
	hasRecentData := false
	var mostRecentPath string
	var mostRecentTime time.Time

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && isDateFolder(info.Name()) {
			recentData, filePath, err := checkDateFolder(path, cutoffTime, fileExtension)
			if err != nil {
				return err
			}
			if !recentData {
				return nil
			}
			hasRecentData = true
			fileTime := info.ModTime()
			if fileTime.After(mostRecentTime) {
				mostRecentTime = fileTime
				mostRecentPath = filePath
			}
		}
		return nil
	})

	if err != nil {
		return false, "", err
	}
	return hasRecentData, mostRecentPath, nil
}

func checkDateFolder(path string, cutoffTime time.Time, fileExtension string) (bool, string, error) {
	folderDate, err := time.Parse(dateFormat, filepath.Base(path))
	if err != nil {
		return false, "", err
	}
	if folderDate.After(cutoffTime) {
		return checkForFiles(path, fileExtension)
	}
	return false, "", nil
}

func checkForFiles(path string, fileExtension string) (bool, string, error) {
	var mostRecentPath string
	var mostRecentTime time.Time

	err := filepath.Walk(path, func(subPath string, subInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !subInfo.IsDir() && filepath.Ext(subInfo.Name()) == fileExtension {
			fileTime := subInfo.ModTime()
			if mostRecentPath == "" || fileTime.After(mostRecentTime) {
				mostRecentTime = fileTime
				mostRecentPath = subPath
			}
		}
		return nil
	})

	return mostRecentPath != "", mostRecentPath, err
}

func isDateFolder(name string) bool {
	_, err := time.Parse(dateFormat, name)
	return err == nil
}

func generateFilePath(baseDir string, now time.Time, folder string) (string, error) {
	dateFolder := now.Format(dateFormat)
	dir := filepath.Join(baseDir, dateFolder, folder)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory structure: %w", err)
	}

	return dir, nil
}

func getUniqueFilename(typeOfTest string, dir string, now time.Time) string {
	base := fmt.Sprintf(typeOfTest+"_test_%s", now.Format("20060102_150405"))
	filename := base + ".json"

	for i := 1; ; i++ {
		if _, err := os.Stat(filepath.Join(dir, filename)); os.IsNotExist(err) {
			return filename
		}
		filename = fmt.Sprintf("%s_%d.json", base, i)
	}
}
