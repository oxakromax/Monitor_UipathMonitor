# Guía de Instalación para el Servicio de Monitoreo

## Antes de Empezar

Antes de ponernos manos a la obra, asegúrate de tener listo y corriendo el servidor API REST del proyecto. Si aún no lo has hecho, puedes seguir la guía en el [repositorio correspondiente](https://github.com/oxakromax/Backend_UipathMonitor).

Este servicio de monitoreo es una pieza clave en nuestra arquitectura de microservicios, interactuando directamente con la API central y la API de UiPath. Aquí vamos a configurarlo para que todo funcione como un reloj.

## Clonando el Repositorio

1. Abre tu terminal y ejecuta los siguientes comandos para clonar el repositorio y moverte a la carpeta del proyecto:
   ```bash
   git clone https://github.com/oxakromax/Monitor_UipathMonitor.git
   cd Monitor_UipathMonitor
   ```

## Configurando las Variables de Entorno

2. Crea un archivo `.env` en la carpeta principal del proyecto. Este archivo debe contener las siguientes variables de entorno:

   ```env
   API_URL=http://127.0.0.1:8080
   DB_KEY=22d89667b85011d157292c03580db4c2
   MONITOR_PASS=234sdfds1
   MONITOR_USER=monitor@localhost.com
   ```

   **Importante**: Asegúrate de usar la misma `DB_KEY`, además las credenciales `MONITOR_*` que configuraste en el servidor y de apuntar correctamente la `API_URL`. 

## Instalación y Ejecución

3. Ahora que ya tenemos todo configurado, es hora de poner en marcha el servicio de monitoreo. Si estás usando Go directamente, puedes hacerlo con los siguientes comandos:

   ```bash
   go mod tidy
   go run main.go
   ```

   Si todo está configurado correctamente, deberías ver un mensaje indicando que el servicio de monitoreo está en ejecución.

## Usando Docker

Si prefieres usar Docker, aquí te dejo un ejemplo de `Dockerfile` que puedes usar para contenerizar tu aplicación:

```Dockerfile
# Usamos la imagen oficial de Go como base
FROM golang:1.21 as builder

# Establecemos el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiamos los archivos del proyecto al contenedor
COPY . .

# Instalamos las dependencias
RUN go mod tidy

# Compilamos la aplicación
RUN go build -o monitor

# Creamos una imagen ligera usando alpine
FROM alpine:latest

# Copiamos el binario compilado desde la imagen builder
COPY --from=builder /app/monitor /app/

# Establecemos el comando por defecto para ejecutar la aplicación
CMD ["/app/monitor"]
```

Para construir y correr la imagen de Docker, puedes usar los siguientes comandos:

```bash
docker build -t monitor_uipathmonitor .
docker run --env-file .env monitor_uipathmonitor
```

## ¿Y Ahora?

¡Listo! Tu servicio de monitoreo debería estar corriendo y comunicándose con la API central. Recuerda mantener un ojo en los logs para asegurarte de que todo esté funcionando como debe.

¡Buena suerte y a monitorear!