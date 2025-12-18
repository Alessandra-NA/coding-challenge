const statisticsService = require('../services/statisticsService');
const { validateMatrix, validateMatrixSize } = require('../utils/validator');
const logger = require('../utils/logger');

/**
 * Controlador para procesar matriz y calcular estadísticas
 * POST /process
 */
async function processMatrix(req, res) {
  try {
    const { original_matrix, rotated_matrix, request_id } = req.body;
    const requestId = request_id || req.requestId;

    // Validar cuerpo de solicitud
    if (!original_matrix) {
      return res.status(400).json({
        error: 'original_matrix es requerido',
        request_id: requestId,
      });
    }

    if (!rotated_matrix) {
      return res.status(400).json({
        error: 'rotated_matrix es requerido',
        request_id: requestId,
      });
    }

    // Validar estructura de matrices
    const originalValidation = validateMatrix(original_matrix);
    if (!originalValidation.valid) {
      return res.status(400).json({
        error: `Matriz original: ${originalValidation.error}`,
        request_id: requestId,
      });
    }

    const rotatedValidation = validateMatrix(rotated_matrix);
    if (!rotatedValidation.valid) {
      return res.status(400).json({
        error: `Matriz rotada: ${rotatedValidation.error}`,
        request_id: requestId,
      });
    }

    // Validar tamaño de matrices
    const originalSizeValidation = validateMatrixSize(original_matrix, 1000, 1000);
    if (!originalSizeValidation.valid) {
      return res.status(400).json({
        error: `Matriz original: ${originalSizeValidation.error}`,
        request_id: requestId,
      });
    }

    const rotatedSizeValidation = validateMatrixSize(rotated_matrix, 1000, 1000);
    if (!rotatedSizeValidation.valid) {
      return res.status(400).json({
        error: `Matriz rotada: ${rotatedSizeValidation.error}`,
        request_id: requestId,
      });
    }

    // Calcular estadísticas (incluye verificación de diagonal para ambas matrices)
    const statistics = statisticsService.calculateStatistics(original_matrix, rotated_matrix);

    // Retornar respuesta
    const response = {
      rotated_matrix: rotated_matrix,
      statistics: statistics,
      request_id: requestId,
    };

    logger.info('Matriz procesada exitosamente', {
      request_id: requestId,
      matrix_size: `${rotated_matrix.length}x${rotated_matrix[0].length}`,
    });

    return res.status(200).json(response);
  } catch (error) {
    logger.error('Error al procesar matriz', {
      error: error.message,
      stack: error.stack,
      request_id: req.requestId,
    });

    return res.status(500).json({
      error: 'error interno del servidor',
      request_id: req.requestId || 'desconocido',
    });
  }
}

/**
 * Endpoint de health check
 * GET /health
 */
function healthCheck(req, res) {
  return res.status(200).json({
    status: 'saludable',
    service: 'node-api',
  });
}

module.exports = {
  processMatrix,
  healthCheck,
};

