const { validateMatrix, validateMatrixSize } = require('../utils/validator');

describe('Validator', () => {
  describe('validateMatrix', () => {
    test('should validate correct matrix', () => {
      const matrix = [[1, 2, 3], [4, 5, 6]];
      const result = validateMatrix(matrix);
      expect(result.valid).toBe(true);
    });

    test('should reject non-array', () => {
      const result = validateMatrix('not an array');
      expect(result.valid).toBe(false);
      expect(result.error).toBe('La matriz debe ser un array');
    });

    test('should reject empty matrix', () => {
      const result = validateMatrix([]);
      expect(result.valid).toBe(false);
      expect(result.error).toBe('La matriz no puede estar vacía');
    });

    test('should reject matrix with empty first row', () => {
      const matrix = [[]];
      const result = validateMatrix(matrix);
      expect(result.valid).toBe(false);
      expect(result.error).toBe('Las filas de la matriz no pueden estar vacías');
    });

    test('should reject non-rectangular matrix', () => {
      const matrix = [[1, 2, 3], [4, 5]];
      const result = validateMatrix(matrix);
      expect(result.valid).toBe(false);
      expect(result.error).toContain('La matriz debe ser rectangular');
      expect(result.error).toContain('la fila 1 tiene longitud 2, se esperaba 3');
    });

    test('should reject matrix with non-array row', () => {
      const matrix = [[1, 2], 'not an array'];
      const result = validateMatrix(matrix);
      expect(result.valid).toBe(false);
      expect(result.error).toContain('La fila 1 no es un array');
    });

    test('should reject matrix with non-number elements', () => {
      const matrix = [[1, 2, 3], [4, '5', 6]];
      const result = validateMatrix(matrix);
      expect(result.valid).toBe(false);
      expect(result.error).toContain('no es un número válido');
      expect(result.error).toContain('[1][1]');
    });

    test('should reject matrix with NaN elements', () => {
      const matrix = [[1, 2, 3], [4, NaN, 6]];
      const result = validateMatrix(matrix);
      expect(result.valid).toBe(false);
      expect(result.error).toContain('no es un número válido');
    });
  });

  describe('validateMatrixSize', () => {
    test('should validate matrix within limits', () => {
      const matrix = [[1, 2], [3, 4]];
      const result = validateMatrixSize(matrix, 10, 10);
      expect(result.valid).toBe(true);
    });

    test('should reject matrix exceeding row limit', () => {
      const matrix = new Array(11).fill([1, 2]);
      const result = validateMatrixSize(matrix, 10, 10);
      expect(result.valid).toBe(false);
      expect(result.error).toContain('La matriz tiene demasiadas filas');
      expect(result.error).toContain('11');
      expect(result.error).toContain('máximo: 10');
    });

    test('should reject matrix exceeding column limit', () => {
      const matrix = [[...new Array(11).fill(1)]];
      const result = validateMatrixSize(matrix, 10, 10);
      expect(result.valid).toBe(false);
      expect(result.error).toContain('La matriz tiene demasiadas columnas');
      expect(result.error).toContain('11');
      expect(result.error).toContain('máximo: 10');
    });
  });
});

