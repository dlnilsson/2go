package togo

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConvertToGoStructs(t *testing.T) {
	tests := []struct {
		name       string
		file       string
		formatType string
		expected   string
	}{
		{
			name:       "Simple JSON",
			file:       "testdata/array.json",
			formatType: "json",
			expected:   "testdata/array.result",
		},
		{
			name:       "complex JSON",
			file:       "testdata/complex.json",
			formatType: "json",
			expected:   "testdata/complex.result",
		},
		{
			name:       "shared JSON",
			file:       "testdata/shared.json",
			formatType: "json",
			expected:   "testdata/shared.result",
		},
		{
			name:       "httpbin modified special char",
			file:       "testdata/httpbin.json",
			formatType: "json",
			expected:   "testdata/httpbin.result",
		},
		{
			name:       "simple YAML",
			file:       "testdata/simple.yaml",
			formatType: "yaml",
			expected:   "testdata/simple.result",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read the file
			data, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", tt.file, err)
			}

			var input interface{}
			if tt.formatType == "json" {
				err = json.Unmarshal(data, &input)
			} else {
				err = yaml.Unmarshal(data, &input)
			}
			if err != nil {
				t.Fatalf("Failed to unmarshal %s: %v", tt.formatType, err)
			}

			goCode, err := ConvertToGoStructs(input, false, tt.formatType)
			if err != nil {
				t.Fatalf("Failed to generate Go code: %v", err)
			}
			expected, err := os.ReadFile(tt.expected)
			if err != nil {
				t.Fatalf("Failed to read expected result file %s: %v", tt.expected, err)
			}
			want := string(expected)
			if goCode != want {
				t.Errorf("Expected:\n%q\nGot:\n%q", want, goCode)
			}
		})
	}
}

func TestStartsWithDigit(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"1", true},
		{"9", true},
		{"A", false},
		{"0", true},
		{"!0", false},
		{"#", false},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%x", i), func(t *testing.T) {
			if got := startsWithDigit(tt.input); got != tt.want {
				t.Errorf("startsWithDigit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fName(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "basic",
			arg:  "basicId",
			want: "BasicID",
		},
		{
			name: "LastInsertId",
			arg:  "LastInsertId",
			want: "LastInsertId",
		},
		{
			name: "kWh",
			arg:  "kWh",
			want: "kWh",
		},
		{
			name: "id",
			arg:  "id",
			want: "ID",
		},
		{
			name: "something",
			arg:  "something",
			want: "Something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fName(tt.arg); got != tt.want {
				t.Errorf("fName() = %v, want %v", got, tt.want)
			}
		})
	}
}
