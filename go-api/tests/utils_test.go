package tests

import (
	"interseguro-challenge/go-api/utils"
	"strings"
	"testing"
)

func TestValidateMatrix(t *testing.T) {
	tests := []struct {
		name       string
		matrix     [][]int
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "valid 3x3 matrix",
			matrix:     [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "valid 2x3 matrix",
			matrix:     [][]int{{1, 2, 3}, {4, 5, 6}},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "empty matrix",
			matrix:     [][]int{},
			wantErr:    true,
			wantErrMsg: "La matriz no puede estar vacía",
		},
		{
			name:       "non-rectangular matrix",
			matrix:     [][]int{{1, 2, 3}, {4, 5}},
			wantErr:    true,
			wantErrMsg: "La matriz debe ser rectangular",
		},
		{
			name:       "matrix with empty row",
			matrix:     [][]int{{1, 2}, {}},
			wantErr:    true,
			wantErrMsg: "La matriz debe ser rectangular",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateMatrix(tt.matrix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMatrix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.wantErrMsg != "" && !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("ValidateMatrix() error message = %q, want to contain %q", err.Error(), tt.wantErrMsg)
				}
			}
		})
	}
}

func TestValidateMatrixSize(t *testing.T) {
	tests := []struct {
		name       string
		matrix     [][]int
		maxRows    int
		maxCols    int
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "matrix within limits",
			matrix:     [][]int{{1, 2, 3}, {4, 5, 6}},
			maxRows:    10,
			maxCols:    10,
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "matrix exceeds row limit",
			matrix:     make([][]int, 11),
			maxRows:    10,
			maxCols:    10,
			wantErr:    true,
			wantErrMsg: "La matriz tiene demasiadas filas",
		},
		{
			name:       "matrix exceeds col limit",
			matrix:     [][]int{make([]int, 11)},
			maxRows:    10,
			maxCols:    10,
			wantErr:    true,
			wantErrMsg: "La matriz tiene demasiadas columnas",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateMatrixSize(tt.matrix, tt.maxRows, tt.maxCols)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMatrixSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.wantErrMsg != "" && !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("ValidateMatrixSize() error message = %q, want to contain %q", err.Error(), tt.wantErrMsg)
				}
			}
		})
	}
}

