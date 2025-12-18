package utils

import (
	"fmt"
)

// ValidateMatrix valida que una matriz sea válida (rectangular, no vacía, contiene solo enteros)
// Retorna error si la matriz es inválida, nil en caso contrario
func ValidateMatrix(matrix [][]int) error {
	if len(matrix) == 0 {
		return fmt.Errorf("La matriz no puede estar vacía")
	}

	if len(matrix[0]) == 0 {
		return fmt.Errorf("Las filas de la matriz no pueden estar vacías")
	}

	// Verificar que todas las filas tengan la misma longitud (matriz rectangular)
	rowLength := len(matrix[0])
	for i, row := range matrix {
		if len(row) != rowLength {
			return fmt.Errorf("La matriz debe ser rectangular: La fila %d tiene longitud %d, se esperaba %d", i, len(row), rowLength)
		}
	}

	return nil
}

// ValidateMatrixSize valida que las dimensiones de la matriz estén dentro de límites aceptables
func ValidateMatrixSize(matrix [][]int, maxRows, maxCols int) error {
	if len(matrix) > maxRows {
		return fmt.Errorf("La matriz tiene demasiadas filas: %d (máximo: %d)", len(matrix), maxRows)
	}

	if len(matrix) > 0 && len(matrix[0]) > maxCols {
		return fmt.Errorf("La matriz tiene demasiadas columnas: %d (máximo: %d)", len(matrix[0]), maxCols)
	}

	return nil
}

