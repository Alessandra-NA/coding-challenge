package services

// RotateMatrix90Clockwise rota una matriz 90 grados en sentido horario
// Para una matriz de tamaño MxN, la matriz rotada será NxM
// Algoritmo: Para la matriz original matrix[i][j], la posición rotada es [j][M-1-i]
// Complejidad temporal: O(M*N), Complejidad espacial: O(M*N)
func RotateMatrix90Clockwise(matrix [][]int) [][]int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return matrix
	}

	rows := len(matrix)
	cols := len(matrix[0])

	// Crear nueva matriz con dimensiones intercambiadas
	rotated := make([][]int, cols)
	for i := range rotated {
		rotated[i] = make([]int, rows)
	}

	// Rotar: nueva[i][j] = original[rows-1-j][i]
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			rotated[j][rows-1-i] = matrix[i][j]
		}
	}

	return rotated
}
