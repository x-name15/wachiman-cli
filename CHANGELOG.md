# Changelog

Todos los cambios notables de este proyecto se documentan aquí.

---
## [1.0.0] - 2026-06-08 — "Ahora si, ya estoy listo pa chambear"

### Añadido
  - Comando `config` y paquete interno `config` para gestionar preferencias locales:
  - Subcomandos: `config get`, `config set [key] [value]`, `config reset`.
  - Opciones de configuración: `watch_interval`, `default_tail`, `output_format`.

### Cambios / Mejoras
  - `watch`: ahora respeta `watch_interval` desde la configuración si `--interval` no fue pasado.
  - `ps`: añadido flag `--running` y `--stopped` para filtrar, y `--output json` para salida en formato JSON.
  - `stats`: añadido `--output json` y formateo/colorado de porcentajes de CPU y memoria.
  - `logs`: cuando no se especifica `--tail`, ahora usa `default_tail` desde la configuración.
  - `main`: registro del `ConfigCmd` en el comando raíz y mejoras en el banner/ayuda.

  ### Archivos modificados
  - `cmd/stats.go` — soporte `--output json`, coloreo de porcentajes.
  - `cmd/ps.go` — filtros `--running`/`--stopped`, `--output json` y salida tabulada.
  - `cmd/logs.go` — uso de `default_tail` desde la configuración cuando `--tail` no fue provisto.
  - `cmd/config.go` — nuevo comando para gestionar la configuración del usuario.
  - `config/config.go` — nuevo paquete para carga/guardado de config en `~/.wachiman/config.json`.
  - `main.go` — registro de `ConfigCmd` y banner/ayuda mejorada.


## [0.2.0] - 2026-06-08 — "Seño, pan con palta y su quinua con manzana"

### Añadido

- `wachiman watch` ahora muestra sparklines de tendencia con caracteres Unicode `▁▂▃▄▅▆▇█`
  - Historial de los últimos 10 ticks por contenedor en memoria
  - Color según el último valor: verde < 50%, amarillo 50–80%, rojo > 80%
- Header dinámico en `wachiman watch` con conteo de contenedores corriendo vs parados
- Fix de pantalla en Windows — `cls` en vez de escape codes ANSI para limpiar correctamente

## [0.1.0] - 2026-06-08 — "Oe wachiman, apura p"

Primera release de Wachiman CLI. Mi causa ha despertado.

### Añadido

- `wachiman ps` — lista todos los contenedores con ID, nombre, imagen, estado y puertos
  - Flag `--running` para filtrar solo contenedores activos
  - Flag `--stopped` para filtrar solo contenedores parados
  - Colores: verde para activos, rojo para parados

- `wachiman stats` — muestra CPU y memoria de los contenedores corriendo
  - Barras de progreso con caracteres `█░` proporcionales al uso
  - Colores según umbral: verde < 50%, amarillo 50–80%, rojo > 80%

- `wachiman watch` — monitor en tiempo real
  - Refresco automático cada 3 segundos (configurable con `--interval` / `-i`)
  - Header dinámico con conteo de contenedores corriendo vs parados
  - Barras de CPU y memoria en vivo
  - Salida limpia con `Ctrl+C`

- `wachiman logs` — muestra los logs de un contenedor
  - Flag `--tail` / `-t` para controlar cuántas líneas mostrar (por defecto: 50)

- `wachiman inspect` — detalles completos de un contenedor
  - IP, puertos expuestos, volúmenes montados y variables de entorno

- `wachiman top` — procesos corriendo dentro de un contenedor

- `wachiman start` — arranca un contenedor parado

- `wachiman stop` — para un contenedor en ejecución

- `wachiman restart` — reinicia un contenedor

- `wachiman prune` — limpia contenedores parados, imágenes sin usar y volúmenes huérfanos
  - Confirmación interactiva antes de borrar
  - Flag `--force` / `-f` para saltar la confirmación
  - Resumen de espacio liberado al finalizar

- Banner ASCII al ejecutar `wachiman` sin argumentos
- Colores en toda la interfaz via `fatih/color`
- Output tabulado y alineado via `text/tabwriter`

