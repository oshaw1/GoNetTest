package dataManagement

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Name    string
	Age     int
	IsValid bool
	Time    time.Time
}

type NestedStruct struct {
	ID      int
	Details TestStruct
}

func TestReadJSONFile(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "test-json-*")
	require.NoError(t, err)
	defer os.RemoveAll(baseDir)

	repo := &Repository{
		baseDir: baseDir,
	}

	tests := []struct {
		name     string
		content  interface{}
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name: "valid simple json",
			content: map[string]interface{}{
				"name": "test",
				"age":  30,
			},
			validate: func(t *testing.T, result interface{}) {
				m, ok := result.(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "test", m["name"])
				assert.Equal(t, float64(30), m["age"]) // JSON numbers are float64
			},
		},
		{
			name: "valid nested json",
			content: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "test",
					"age":  30,
				},
			},
			validate: func(t *testing.T, result interface{}) {
				m, ok := result.(map[string]interface{})
				require.True(t, ok)
				user, ok := m["user"].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "test", user["name"])
			},
		},
		{
			name:    "invalid json",
			content: "invalid{json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			tempFile, err := os.CreateTemp(baseDir, "test-*.json")
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			// Write test content
			var content []byte
			if str, ok := tt.content.(string); ok {
				content = []byte(str)
			} else {
				content, err = json.Marshal(tt.content)
				require.NoError(t, err)
			}
			err = os.WriteFile(tempFile.Name(), content, 0644)
			require.NoError(t, err)

			// Test reading
			result, err := repo.readJSONFile(tempFile.Name())
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestConvertToMap(t *testing.T) {
	repo := &Repository{}
	now := time.Now()

	tests := []struct {
		name     string
		input    interface{}
		wantErr  bool
		validate func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "struct conversion",
			input: TestStruct{
				Name:    "test",
				Age:     30,
				IsValid: true,
				Time:    now,
			},
			validate: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "test", result["Name"])
				assert.Equal(t, 30, result["Age"])
				assert.Equal(t, true, result["IsValid"])
				assert.Equal(t, now, result["Time"])
			},
		},
		{
			name: "pointer to struct",
			input: &TestStruct{
				Name:    "test",
				Age:     30,
				IsValid: true,
			},
			validate: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "test", result["Name"])
				assert.Equal(t, 30, result["Age"])
				assert.Equal(t, true, result["IsValid"])
			},
		},
		{
			name: "nested struct",
			input: NestedStruct{
				ID: 1,
				Details: TestStruct{
					Name: "test",
					Age:  30,
				},
			},
			validate: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, 1, result["ID"])
				details, ok := result["Details"].(TestStruct)
				require.True(t, ok)
				assert.Equal(t, "test", details.Name)
				assert.Equal(t, 30, details.Age)
			},
		},
		{
			name: "existing map",
			input: map[string]interface{}{
				"name": "test",
				"age":  30,
			},
			validate: func(t *testing.T, result map[string]interface{}) {
				assert.Equal(t, "test", result["name"])
				assert.Equal(t, 30, result["age"])
			},
		},
		{
			name:    "unsupported type",
			input:   42,
			wantErr: true,
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.convertToMap(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestIsDateFolder(t *testing.T) {
	tests := []struct {
		name     string
		folder   string
		expected bool
	}{
		{
			name:     "valid date format",
			folder:   "2024-03-25",
			expected: true,
		},
		{
			name:     "invalid date format",
			folder:   "2024/03/25",
			expected: false,
		},
		{
			name:     "not a date",
			folder:   "test-folder",
			expected: false,
		},
		{
			name:     "empty string",
			folder:   "",
			expected: false,
		},
		{
			name:     "partial date",
			folder:   "2024-03",
			expected: false,
		},
		{
			name:     "invalid month",
			folder:   "2024-13-25",
			expected: false,
		},
		{
			name:     "invalid day",
			folder:   "2024-03-32",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDateFolder(tt.folder)
			assert.Equal(t, tt.expected, result)
		})
	}
}
