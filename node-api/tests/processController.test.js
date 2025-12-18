const request = require('supertest');
const app = require('../index');

describe('ProcessController', () => {
  describe('POST /process', () => {
    test('should process matrix and return statistics', async () => {
      const matrix = [
        [1, 2, 3],
        [4, 5, 6],
        [7, 8, 9]
      ];

      const response = await request(app)
        .post('/process')
        .send({
          original_matrix: matrix,
          rotated_matrix: matrix,
          request_id: 'test-123'
        })
        .expect(200);

      expect(response.body).toHaveProperty('rotated_matrix');
      expect(response.body).toHaveProperty('statistics');
      expect(response.body.statistics).toHaveProperty('max_value');
      expect(response.body.statistics).toHaveProperty('min_value');
      expect(response.body.statistics).toHaveProperty('average');
      expect(response.body.statistics).toHaveProperty('total_sum');
      expect(response.body.request_id).toBe('test-123');
    });

    test('should return 400 for missing original_matrix', async () => {
      const response = await request(app)
        .post('/process')
        .send({
          rotated_matrix: [[1, 2], [3, 4]]
        })
        .expect(400);

      expect(response.body).toHaveProperty('error');
      expect(response.body.error).toBe('original_matrix es requerido');
    });

    test('should return 400 for missing rotated_matrix', async () => {
      const response = await request(app)
        .post('/process')
        .send({
          original_matrix: [[1, 2], [3, 4]]
        })
        .expect(400);

      expect(response.body).toHaveProperty('error');
      expect(response.body.error).toBe('rotated_matrix es requerido');
    });

    test('should return 400 for empty original_matrix', async () => {
      const response = await request(app)
        .post('/process')
        .send({
          original_matrix: [],
          rotated_matrix: [[1, 2], [3, 4]]
        })
        .expect(400);

      expect(response.body).toHaveProperty('error');
      expect(response.body.error).toContain('Matriz original');
      expect(response.body.error).toContain('La matriz no puede estar vacía');
    });

    test('should return 400 for empty rotated_matrix', async () => {
      const response = await request(app)
        .post('/process')
        .send({
          original_matrix: [[1, 2], [3, 4]],
          rotated_matrix: []
        })
        .expect(400);

      expect(response.body).toHaveProperty('error');
      expect(response.body.error).toContain('Matriz rotada');
      expect(response.body.error).toContain('La matriz no puede estar vacía');
    });

    test('should return 400 for non-rectangular original_matrix', async () => {
      const response = await request(app)
        .post('/process')
        .send({
          original_matrix: [[1, 2, 3], [4, 5]],
          rotated_matrix: [[1, 2], [3, 4]]
        })
        .expect(400);

      expect(response.body).toHaveProperty('error');
      expect(response.body.error).toContain('Matriz original');
      expect(response.body.error).toContain('La matriz debe ser rectangular');
    });

    test('should return 400 for non-rectangular rotated_matrix', async () => {
      const response = await request(app)
        .post('/process')
        .send({
          original_matrix: [[1, 2], [3, 4]],
          rotated_matrix: [[1, 2], [3]]
        })
        .expect(400);

      expect(response.body).toHaveProperty('error');
      expect(response.body.error).toContain('Matriz rotada');
      expect(response.body.error).toContain('La matriz debe ser rectangular');
    });

    test('should return 400 for original_matrix exceeding size limit', async () => {
      const largeMatrix = new Array(1001).fill([1, 2]);
      const response = await request(app)
        .post('/process')
        .send({
          original_matrix: largeMatrix,
          rotated_matrix: [[1, 2], [3, 4]]
        })
        .expect(400);

      expect(response.body).toHaveProperty('error');
      expect(response.body.error).toContain('Matriz original');
      expect(response.body.error).toContain('La matriz tiene demasiadas filas');
    });

    test('should return 400 for rotated_matrix exceeding size limit', async () => {
      const largeMatrix = [[...new Array(1001).fill(1)]];
      const response = await request(app)
        .post('/process')
        .send({
          original_matrix: [[1, 2], [3, 4]],
          rotated_matrix: largeMatrix
        })
        .expect(400);

      expect(response.body).toHaveProperty('error');
      expect(response.body.error).toContain('Matriz rotada');
      expect(response.body.error).toContain('La matriz tiene demasiadas columnas');
    });
  });

  describe('GET /health', () => {
    test('should return healthy status in Spanish', async () => {
      const response = await request(app)
        .get('/health')
        .expect(200);

      expect(response.body).toHaveProperty('status');
      expect(response.body.status).toBe('saludable');
      expect(response.body).toHaveProperty('service');
      expect(response.body.service).toBe('node-api');
    });
  });
});

