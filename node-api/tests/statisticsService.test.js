const { calculateStatistics, checkIfDiagonal } = require('../services/statisticsService');

describe('StatisticsService', () => {
  describe('calculateStatistics', () => {
    test('should calculate correct statistics for a 3x3 matrix', () => {
      const originalMatrix = [
        [1, 2, 3],
        [4, 5, 6],
        [7, 8, 9]
      ];
      const rotatedMatrix = [
        [7, 4, 1],
        [8, 5, 2],
        [9, 6, 3]
      ];

      const result = calculateStatistics(originalMatrix, rotatedMatrix);

      expect(result.max_value).toBe(9);
      expect(result.min_value).toBe(1);
      expect(result.average).toBe(5);
      expect(result.total_sum).toBe(45);
      expect(result.rotated_is_diagonal).toBe(false);
      expect(result.calculation_time_ms).toBeGreaterThanOrEqual(0);
    });

    test('should calculate correct statistics for a 2x2 matrix', () => {
      const originalMatrix = [
        [10, 20],
        [30, 40]
      ];
      const rotatedMatrix = [
        [30, 10],
        [40, 20]
      ];

      const result = calculateStatistics(originalMatrix, rotatedMatrix);

      expect(result.max_value).toBe(40);
      expect(result.min_value).toBe(10);
      expect(result.average).toBe(25);
      expect(result.total_sum).toBe(100);
      expect(result.rotated_is_diagonal).toBe(false);
    });

    test('should identify diagonal matrix', () => {
      const originalMatrix = [
        [1, 0, 0],
        [0, 2, 0],
        [0, 0, 3]
      ];
      const rotatedMatrix = [
        [0, 0, 1],
        [0, 2, 0],
        [3, 0, 0]
      ];

      const result = calculateStatistics(originalMatrix, rotatedMatrix);

      expect(result.original_is_diagonal).toBe(true);
      expect(result.rotated_is_diagonal).toBe(false);
      expect(result.total_sum).toBe(6);
    });

    test('should handle matrix with negative numbers', () => {
      const originalMatrix = [
        [-5, -2],
        [-1, 0]
      ];
      const rotatedMatrix = [
        [-1, -5],
        [0, -2]
      ];

      const result = calculateStatistics(originalMatrix, rotatedMatrix);

      expect(result.max_value).toBe(0);
      expect(result.min_value).toBe(-5);
      expect(result.average).toBe(-2);
      expect(result.total_sum).toBe(-8);
    });

    test('should handle single element matrix', () => {
      const originalMatrix = [[42]];
      const rotatedMatrix = [[42]];

      const result = calculateStatistics(originalMatrix, rotatedMatrix);

      expect(result.max_value).toBe(42);
      expect(result.min_value).toBe(42);
      expect(result.average).toBe(42);
      expect(result.total_sum).toBe(42);
      expect(result.original_is_diagonal).toBe(true);
      expect(result.rotated_is_diagonal).toBe(true);
    });

    test('should throw error for empty rotated matrix', () => {
      expect(() => {
        calculateStatistics([[1, 2], [3, 4]], []);
      }).toThrow('La matriz rotada no puede estar vacía');
    });

    test('should handle rectangular matrix (2x3)', () => {
      const originalMatrix = [
        [1, 2, 3],
        [4, 5, 6]
      ];
      const rotatedMatrix = [
        [4, 1],
        [5, 2],
        [6, 3]
      ];

      const result = calculateStatistics(originalMatrix, rotatedMatrix);

      expect(result.max_value).toBe(6);
      expect(result.min_value).toBe(1);
      expect(result.average).toBe(3.5);
      expect(result.total_sum).toBe(21);
      expect(result.rotated_is_diagonal).toBe(false);
    });
  });

  describe('checkIfDiagonal', () => {
    test('should return true for diagonal matrix', () => {
      const matrix = [
        [5, 0, 0],
        [0, 3, 0],
        [0, 0, 7]
      ];

      expect(checkIfDiagonal(matrix)).toBe(true);
    });

    test('should return false for non-diagonal square matrix', () => {
      const matrix = [
        [1, 2, 3],
        [4, 5, 6],
        [7, 8, 9]
      ];

      expect(checkIfDiagonal(matrix)).toBe(false);
    });

    test('should return false for rectangular matrix', () => {
      const matrix = [
        [1, 2, 3],
        [4, 5, 6]
      ];

      expect(checkIfDiagonal(matrix)).toBe(false);
    });

    test('should return true for single element matrix', () => {
      expect(checkIfDiagonal([[42]])).toBe(true);
    });

    test('should return false for matrix with non-zero off-diagonal', () => {
      const matrix = [
        [1, 0],
        [0, 1]
      ];
      expect(checkIfDiagonal(matrix)).toBe(true);

      const matrix2 = [
        [1, 1],
        [0, 1]
      ];
      expect(checkIfDiagonal(matrix2)).toBe(false);
    });
  });
});

