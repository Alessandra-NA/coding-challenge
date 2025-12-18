const logger = require('../utils/logger');

/**
 * Calcula estadísticas para una matriz rotada y verifica diagonal para ambas matrices
 * @param {number[][]} originalMatrix - La matriz original
 * @param {number[][]} rotatedMatrix - La matriz rotada a analizar
 * @returns {Object} Objeto de estadísticas con max, min, average, totalSum, original_is_diagonal y rotated_is_diagonal
 */
function calculateStatistics(originalMatrix, rotatedMatrix) {
  const startTime = Date.now();

  if (!rotatedMatrix || rotatedMatrix.length === 0) {
    throw new Error('La matriz rotada no puede estar vacía');
  }

  let maxValue = rotatedMatrix[0][0];
  let minValue = rotatedMatrix[0][0];
  let totalSum = 0;
  let elementCount = 0;

  // Un solo recorrido por la matriz rotada para calcular todas las estadísticas (complejidad O(n))
  for (let i = 0; i < rotatedMatrix.length; i++) {
    for (let j = 0; j < rotatedMatrix[i].length; j++) {
      const value = rotatedMatrix[i][j];
      totalSum += value;
      elementCount++;

      if (value > maxValue) {
        maxValue = value;
      }
      if (value < minValue) {
        minValue = value;
      }
    }
  }

  const average = elementCount > 0 ? totalSum / elementCount : 0;
  
  // Verificar si ambas matrices son diagonales
  const originalIsDiagonal = checkIfDiagonal(originalMatrix);
  const rotatedIsDiagonal = checkIfDiagonal(rotatedMatrix);

  const calculationTime = Date.now() - startTime;

  const statistics = {
    max_value: maxValue,
    min_value: minValue,
    average: parseFloat(average.toFixed(6)), // Redondear a 6 decimales
    total_sum: totalSum,
    original_is_diagonal: originalIsDiagonal,
    rotated_is_diagonal: rotatedIsDiagonal,
    calculation_time_ms: calculationTime,
  };

  logger.debug('Estadísticas calculadas', { statistics });
  return statistics;
}

/**
 * Verifica si una matriz es diagonal (solo para matrices cuadradas)
 * Una matriz diagonal tiene elementos no cero solo en la diagonal principal
 * @param {number[][]} matrix - La matriz a verificar
 * @returns {boolean} True si la matriz es diagonal, false en caso contrario
 */
function checkIfDiagonal(matrix) {
  // Solo las matrices cuadradas pueden ser diagonales
  if (matrix.length === 0 || matrix.length !== matrix[0].length) {
    return false;
  }

  const n = matrix.length;

  // Verificar todos los elementos: los elementos de la diagonal pueden tener cualquier valor,
  // pero los elementos no diagonales deben ser cero
  for (let i = 0; i < n; i++) {
    for (let j = 0; j < n; j++) {
      // Si no está en la diagonal principal y el valor no es cero
      if (i !== j && matrix[i][j] !== 0) {
        return false;
      }
    }
  }

  return true;
}

module.exports = {
  calculateStatistics,
  checkIfDiagonal,
};

