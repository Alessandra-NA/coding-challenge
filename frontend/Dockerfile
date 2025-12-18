# Dockerfile para frontend con configuración dinámica
FROM nginx:alpine

# sed ya está incluido en alpine, no necesitamos instalar nada adicional

# Copiar archivos estáticos
COPY index.html /usr/share/nginx/html/
COPY app.js /usr/share/nginx/html/
COPY style.css /usr/share/nginx/html/

# Copiar configuración de nginx template
COPY nginx.conf.template /etc/nginx/nginx.conf.template

# Copiar script de entrada
COPY entrypoint.sh /entrypoint.sh

# Convertir finales de línea a Unix (LF) y dar permisos de ejecución
RUN sed -i 's/\r$//' /entrypoint.sh && chmod +x /entrypoint.sh

# Exponer puerto (Cloud Run inyecta PORT)
EXPOSE 80

# Sobrescribir el ENTRYPOINT de nginx con nuestro script
ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]
