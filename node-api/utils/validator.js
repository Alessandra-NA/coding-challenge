/**
 * Valida que una matriz sea válida (rectangular, no vacía)
 * @param {number[][]} matrix - La matriz a validar
 * @returns {Object} - { valid: boolean, error: string }
 */
function validateMatrix(matrix) {
  if (!Array.isArray(matrix)) {
    return { valid: false, error: 'La matriz debe ser un array' };
  }

  if (matrix.length === 0) {
    return { valid: false, error: 'La matriz no puede estar vacía' };
  }

  if (!Array.isArray(matrix[0]) || matrix[0].length === 0) {
    return { valid: false, error: 'Las filas de la matriz no pueden estar vacías' };
  }

  // Verificar que todas las filas tengan la misma longitud (matriz rectangular)
  const rowLength = matrix[0].length;
  for (let i = 0; i < matrix.length; i++) {
    if (!Array.isArray(matrix[i])) {
      return { valid: false, error: `La fila ${i} no es un array` };
    }

    if (matrix[i].length !== rowLength) {
      return { valid: false, error: `La matriz debe ser rectangular: la fila ${i} tiene longitud ${matrix[i].length}, se esperaba ${rowLength}` };
    }

    // Validar que todos los elementos sean números
    for (let j = 0; j < matrix[i].length; j++) {
      if (typeof matrix[i][j] !== 'number' || isNaN(matrix[i][j])) {
        return { valid: false, error: `El elemento de la matriz en [${i}][${j}] no es un número válido` };
      }
    }
  }

  return { valid: true };
}

/**
 * Valida el tamaño de la matriz dentro de límites
 * @param {number[][]} matrix - La matriz a validar
 * @param {number} maxRows - Número máximo de filas permitidas
 * @param {number} maxCols - Número máximo de columnas permitidas
 * @returns {Object} - { valid: boolean, error: string }
 */
function validateMatrixSize(matrix, maxRows = 1000, maxCols = 1000) {
  if (matrix.length > maxRows) {
    return { valid: false, error: `La matriz tiene demasiadas filas: ${matrix.length} (máximo: ${maxRows})` };
  }

  if (matrix.length > 0 && matrix[0].length > maxCols) {
    return { valid: false, error: `La matriz tiene demasiadas columnas: ${matrix[0].length} (máximo: ${maxCols})` };
  }

  return { valid: true };
}

module.exports = {
  validateMatrix,
  validateMatrixSize,
};

