const { v4: uuidv4 } = require('uuid');
const logger = require('../utils/logger');

/**
 * Middleware para agregar ID de solicitud y registrar solicitudes
 */
function loggerMiddleware(req, res, next) {
  // Generar o recuperar ID de solicitud
  const requestId = req.headers['x-request-id'] || uuidv4();
  req.requestId = requestId;
  res.setHeader('X-Request-ID', requestId);

  // Registrar tiempo de inicio
  const startTime = Date.now();

  // Registrar inicio de solicitud
  logger.info('Solicitud recibida', {
    request_id: requestId,
    method: req.method,
    path: req.path,
    ip: req.ip,
    user_agent: req.get('user-agent'),
  });

  // Sobrescribir res.json para registrar respuesta
  const originalJson = res.json.bind(res);
  res.json = function (body) {
    const duration = Date.now() - startTime;

    logger.info('Solicitud completada', {
      request_id: requestId,
      method: req.method,
      path: req.path,
      status_code: res.statusCode,
      duration_ms: duration,
    });

    return originalJson(body);
  };

  next();
}

module.exports = loggerMiddleware;

