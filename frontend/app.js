// Configuration
// API_BASE_URL se carga desde config.js (generado dinámicamente)
// Fallback a localhost si config.js no está disponible
const API_BASE_URL = (window.APP_CONFIG && window.APP_CONFIG.API_BASE_URL) || 'http://localhost:8080';

// State
let authToken = localStorage.getItem('authToken');
let tokenExpiry = localStorage.getItem('tokenExpiry');

// DOM Elements
const loginSection = document.getElementById('loginSection');
const mainSection = document.getElementById('mainSection');
const loginForm = document.getElementById('loginForm');
const loginError = document.getElementById('loginError');
const matrixInput = document.getElementById('matrixInput');
const rotateBtn = document.getElementById('rotateBtn');
const rotateError = document.getElementById('rotateError');
const resultsSection = document.getElementById('resultsSection');
let originalMatrixDiv = document.getElementById('originalMatrix');
let rotatedMatrixDiv = document.getElementById('rotatedMatrix');
const statisticsDiv = document.getElementById('statistics');
const processingTimeSpan = document.getElementById('processingTime');
const loadingIndicator = document.getElementById('loadingIndicator');
const copyNotification = document.getElementById('copyNotification');

// Check if user is already logged in
if (authToken && tokenExpiry && Date.now() < parseInt(tokenExpiry)) {
    showMainSection();
} else {
    clearAuth();
}

// Variables para almacenar las matrices mostradas
let currentOriginalMatrix = null;
let currentRotatedMatrix = null;

// Login form handler
loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    hideError(loginError);

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    try {
        const response = await fetch(`${API_BASE_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        const data = await response.json();

        if (!response.ok) {
            showError(loginError, data.error || 'Error de autenticación');
            return;
        }

        // Store token
        authToken = data.token;
        const expiresIn = data.expires_in * 1000; // Convert to milliseconds
        tokenExpiry = Date.now() + expiresIn;

        localStorage.setItem('authToken', authToken);
        localStorage.setItem('tokenExpiry', tokenExpiry.toString());

        showMainSection();
    } catch (error) {
        showError(loginError, 'Error de conexión. Verifique que el servidor esté ejecutándose.');
        console.error('Login error:', error);
    }
});

// Rotate button handler
rotateBtn.addEventListener('click', async () => {
    hideError(rotateError);
    hideResults();

    const matrixText = matrixInput.value.trim();

    if (!matrixText) {
        showError(rotateError, 'Por favor ingrese una matriz');
        return;
    }

    // Parse matrix
    let matrix;
    try {
        matrix = JSON.parse(matrixText);
    } catch (error) {
        showError(rotateError, 'Formato JSON inválido. Use el formato: [[1,2,3],[4,5,6]]');
        return;
    }

    // Validate matrix structure
    if (!Array.isArray(matrix) || matrix.length === 0 || !Array.isArray(matrix[0])) {
        showError(rotateError, 'La matriz debe ser un array de arrays de números');
        return;
    }

    // Check if token is still valid
    if (!authToken || Date.now() >= parseInt(tokenExpiry || 0)) {
        showError(rotateError, 'Sesión expirada. Por favor inicie sesión nuevamente.');
        clearAuth();
        showLoginSection();
        return;
    }

    // Show loading
    showLoading();

    try {
        const response = await fetch(`${API_BASE_URL}/rotate`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`,
            },
            body: JSON.stringify({ matrix }),
        });

        const data = await response.json();

        hideLoading();

        if (!response.ok) {
            if (response.status === 401) {
                showError(rotateError, 'Token expirado. Por favor inicie sesión nuevamente.');
                clearAuth();
                showLoginSection();
            } else {
                showError(rotateError, data.error || 'Error al procesar la matriz');
            }
            return;
        }

        // Display results
        displayResults(data);
    } catch (error) {
        hideLoading();
        showError(rotateError, 'Error de conexión. Verifique que el servidor esté ejecutándose.');
        console.error('Rotate error:', error);
    }
});

// Display results
function displayResults(data) {
    // Almacenar matrices para copiar al portapapeles
    currentOriginalMatrix = data.original_matrix;
    currentRotatedMatrix = data.rotated_matrix;

    // Display original matrix
    originalMatrixDiv.innerHTML = renderMatrix(data.original_matrix);
    
    // Display diagonal info for original matrix
    const originalDiagonalInfo = document.getElementById('originalDiagonalInfo');
    if (data.statistics && data.statistics.original_is_diagonal !== undefined) {
        originalDiagonalInfo.innerHTML = `<p class="diagonal-badge ${data.statistics.original_is_diagonal ? 'diagonal-yes' : 'diagonal-no'}">
            ${data.statistics.original_is_diagonal ? '✓ Es diagonal' : '✗ No es diagonal'}
        </p>`;
    } else {
        originalDiagonalInfo.innerHTML = '';
    }

    // Display rotated matrix
    rotatedMatrixDiv.innerHTML = renderMatrix(data.rotated_matrix);
    
    // Display diagonal info for rotated matrix
    const rotatedDiagonalInfo = document.getElementById('rotatedDiagonalInfo');
    if (data.statistics && data.statistics.rotated_is_diagonal !== undefined) {
        rotatedDiagonalInfo.innerHTML = `<p class="diagonal-badge ${data.statistics.rotated_is_diagonal ? 'diagonal-yes' : 'diagonal-no'}">
            ${data.statistics.rotated_is_diagonal ? '✓ Es diagonal' : '✗ No es diagonal'}
        </p>`;
    } else {
        rotatedDiagonalInfo.innerHTML = '';
    }

    // Display statistics
    if (data.statistics) {
        statisticsDiv.innerHTML = renderStatistics(data.statistics);
        processingTimeSpan.textContent = data.processing_time_ms.toFixed(2);
    }

    // Agregar event listeners para copiar al portapapeles
    setupMatrixCopyListeners();

    // Show results section
    resultsSection.classList.remove('hidden');
}

// Render matrix as HTML
function renderMatrix(matrix) {
    return matrix.map(row => {
        const cells = row.map(cell => `<span class="matrix-cell">${cell}</span>`).join('');
        return `<div class="matrix-row">${cells}</div>`;
    }).join('');
}

// Render statistics as HTML
function renderStatistics(stats) {
    return `
        <div class="stat-item">
            <label>Valor Máximo</label>
            <value>${stats.max_value}</value>
        </div>
        <div class="stat-item">
            <label>Valor Mínimo</label>
            <value>${stats.min_value}</value>
        </div>
        <div class="stat-item">
            <label>Promedio</label>
            <value>${stats.average.toFixed(2)}</value>
        </div>
        <div class="stat-item">
            <label>Suma Total</label>
            <value>${stats.total_sum}</value>
        </div>
    `;
}

// UI Helper functions
function showMainSection() {
    loginSection.classList.add('hidden');
    mainSection.classList.remove('hidden');
}

function showLoginSection() {
    loginSection.classList.remove('hidden');
    mainSection.classList.add('hidden');
}

function showError(errorElement, message) {
    errorElement.textContent = message;
    errorElement.classList.add('show');
}

function hideError(errorElement) {
    errorElement.classList.remove('show');
}

function hideResults() {
    resultsSection.classList.add('hidden');
}

function showLoading() {
    loadingIndicator.classList.remove('hidden');
    rotateBtn.disabled = true;
}

function hideLoading() {
    loadingIndicator.classList.add('hidden');
    rotateBtn.disabled = false;
}

function clearAuth() {
    authToken = null;
    tokenExpiry = null;
    localStorage.removeItem('authToken');
    localStorage.removeItem('tokenExpiry');
}

function showCopyNotification() {
    copyNotification.classList.remove('hidden');
    setTimeout(() => {
        copyNotification.classList.add('hidden');
    }, 2000);
}

// Configurar event listeners para copiar matrices al portapapeles
function setupMatrixCopyListeners() {
    // Función helper para copiar al portapapeles
    const copyToClipboard = async (text) => {
        try {
            await navigator.clipboard.writeText(text);
            showCopyNotification();
        } catch (err) {
            // Fallback para navegadores antiguos
            const textArea = document.createElement('textarea');
            textArea.value = text;
            textArea.style.position = 'fixed';
            textArea.style.opacity = '0';
            document.body.appendChild(textArea);
            textArea.select();
            document.execCommand('copy');
            document.body.removeChild(textArea);
            showCopyNotification();
        }
    };

    // Obtener referencias actuales de los divs
    const originalDiv = document.getElementById('originalMatrix');
    const rotatedDiv = document.getElementById('rotatedMatrix');

    // Remover listeners anteriores clonando los nodos (limpia todos los event listeners)
    if (originalDiv) {
        const newOriginal = originalDiv.cloneNode(true);
        originalDiv.parentNode.replaceChild(newOriginal, originalDiv);
        // Actualizar referencia global
        originalMatrixDiv = newOriginal;
        newOriginal.addEventListener('click', async () => {
            if (currentOriginalMatrix) {
                await copyToClipboard(JSON.stringify(currentOriginalMatrix));
            }
        });
    }

    if (rotatedDiv) {
        const newRotated = rotatedDiv.cloneNode(true);
        rotatedDiv.parentNode.replaceChild(newRotated, rotatedDiv);
        // Actualizar referencia global
        rotatedMatrixDiv = newRotated;
        newRotated.addEventListener('click', async () => {
            if (currentRotatedMatrix) {
                await copyToClipboard(JSON.stringify(currentRotatedMatrix));
            }
        });
    }
}

