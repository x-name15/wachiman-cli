# Changelog

Todos los cambios notables de este proyecto se documentan aquí.

---

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
