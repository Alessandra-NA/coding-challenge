#!/bin/sh
# Script de entrada para generar config.js dinámicamente desde variables de entorno

# Obtener API_BASE_URL desde variable de entorno, con fallback
API_BASE_URL=${API_BASE_URL:-http://localhost:8080}

# Obtener PORT desde variable de entorno (Cloud Run lo inyecta)
PORT=${PORT:-80}

# Generar archivo config.js
cat > /usr/share/nginx/html/config.js <<EOF
// Configuración generada dinámicamente en tiempo de ejecución
window.APP_CONFIG = {
    API_BASE_URL: '${API_BASE_URL}'
};
EOF

# Generar configuración de nginx con el puerto correcto usando sed
sed "s/\${PORT}/$PORT/g" /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf

# Verificar que nginx.conf se generó correctamente (para debugging)
echo "Nginx configurado para escuchar en puerto: $PORT"

# Iniciar nginx en foreground
exec nginx -g 'daemon off;'
