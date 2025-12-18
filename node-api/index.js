require('dotenv').config();
const express = require('express');
const cors = require('cors');
const logger = require('./utils/logger');
const loggerMiddleware = require('./middleware/loggerMiddleware');
const processController = require('./controllers/processController');

const app = express();

// Middleware
app.use(cors({
  origin: '*', // En producción, especificar orígenes exactos
  methods: ['GET', 'POST', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization', 'X-Request-ID'],
  exposedHeaders: ['X-Request-ID'],
}));

app.use(express.json({ limit: '10mb' })); // Parsear cuerpos JSON
app.use(loggerMiddleware); // Logging de solicitudes con IDs de solicitud

// Middleware de manejo de errores
app.use((err, req, res, next) => {
  logger.error('Error no manejado', {
    error: err.message,
    stack: err.stack,
    request_id: req.requestId,
  });

  res.status(500).json({
    error: 'error interno del servidor',
    request_id: req.requestId || 'desconocido',
  });
});

// Rutas
app.get('/health', processController.healthCheck);
app.post('/process', processController.processMatrix);

// Manejador 404
app.use((req, res) => {
  res.status(404).json({
    error: 'No encontrado',
    request_id: req.requestId || 'desconocido',
  });
});

// Obtener puerto desde variable de entorno
const port = process.env.PORT || process.env.API2_PORT || 3001;

// Iniciar servidor solo si este archivo se ejecuta directamente (no en tests)
if (require.main === module) {
  app.listen(port, () => {
    logger.info(`Servidor Node.js API iniciado en el puerto ${port}`);
  });
}

module.exports = app;

