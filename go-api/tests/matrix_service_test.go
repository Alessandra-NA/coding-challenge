package tests

import (
	"interseguro-challenge/go-api/services"
	"testing"
)

func TestRotateMatrix90Clockwise(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		expected [][]int
	}{
		{
			name:     "3x3 matrix",
			input:    [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			expected: [][]int{{7, 4, 1}, {8, 5, 2}, {9, 6, 3}},
		},
		{
			name:     "2x3 matrix",
			input:    [][]int{{1, 2, 3}, {4, 5, 6}},
			expected: [][]int{{4, 1}, {5, 2}, {6, 3}},
		},
		{
			name:     "3x2 matrix",
			input:    [][]int{{1, 2}, {3, 4}, {5, 6}},
			expected: [][]int{{5, 3, 1}, {6, 4, 2}},
		},
		{
			name:     "1x1 matrix",
			input:    [][]int{{42}},
			expected: [][]int{{42}},
		},
		{
			name:     "2x2 matrix",
			input:    [][]int{{1, 2}, {3, 4}},
			expected: [][]int{{3, 1}, {4, 2}},
		},
		{
			name:     "empty matrix",
			input:    [][]int{},
			expected: [][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := services.RotateMatrix90Clockwise(tt.input)

			if !matricesEqual(result, tt.expected) {
				t.Errorf("RotateMatrix90Clockwise() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func matricesEqual(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}

	return true
}

func TestRotateMatrix90Clockwise_Dimensions(t *testing.T) {
	tests := []struct {
		name        string
		input       [][]int
		expectedRows int
		expectedCols int
	}{
		{
			name:        "3x3 -> 3x3",
			input:       [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			expectedRows: 3,
			expectedCols: 3,
		},
		{
			name:        "2x3 -> 3x2",
			input:       [][]int{{1, 2, 3}, {4, 5, 6}},
			expectedRows: 3,
			expectedCols: 2,
		},
		{
			name:        "3x2 -> 2x3",
			input:       [][]int{{1, 2}, {3, 4}, {5, 6}},
			expectedRows: 2,
			expectedCols: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := services.RotateMatrix90Clockwise(tt.input)

			if len(result) != tt.expectedRows {
				t.Errorf("Expected %d rows, got %d", tt.expectedRows, len(result))
			}

			if len(result) > 0 && len(result[0]) != tt.expectedCols {
				t.Errorf("Expected %d columns, got %d", tt.expectedCols, len(result[0]))
			}
		})
	}
}

