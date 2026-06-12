# Changelog

Todos los cambios notables de este proyecto se documentan aquí.

---

## [1.9.9] - 2026-06-11 — "Achucha, ahora puedo avisar por radio si encuentro algo"

### Añadido

- `wachiman monitor` ahora soporta notificaciones via webhook
  - Compatible con Slack (incoming webhook) y Discord
  - Notifica cuando un contenedor se cae, cuando es reiniciado, y cuando vuelve a estar activo
  - Notificaciones asíncronas — no bloquean el monitor si el webhook tarda
  - Se activa configurando `webhook_url` y `webhook_type` en la config
- Nuevo paquete interno `internal/webhook` — notifier genérico para Slack y Discord
  - Slack: bloques con sección y contexto
  - Discord: embeds con color según severidad (verde para ok, rojo para caído)
- `wachiman audit` mejorado con tres checks nuevos:
  - **Secrets en variables de entorno** — detecta `PASSWORD`, `TOKEN`, `SECRET`, `API_KEY` y similares en texto plano
  - **Imagen con tag `latest`** — advierte sobre falta de reproducibilidad y riesgos de seguridad
  - **Flag `--fix`** — intenta corregir automáticamente lo que sea posible (límite de memoria por defecto: 512MB)
- `config` ampliado con dos campos nuevos:
  - `webhook_url` — URL del webhook de Slack o Discord
  - `webhook_type` — tipo de webhook (`slack` o `discord`)

### Archivos modificados

- `cmd/audit.go` — checks de secrets, tag latest, y flag `--fix`
- `cmd/monitor.go` — integración de notificaciones via webhook
- `config/config.go` — campos `WebhookURL` y `WebhookType`
- `docker/client.go` — nuevo método `SetMemoryLimit`
- `internal/webhook/webhook.go` — nuevo paquete de notificaciones

---
## [1.9.5] - 2026-06-11 — "Nicagando graba en esa calidad la wbda"

### Añadido

- Nuevo comando `wachiman version` — muestra la versión instalada del CLI
- Workflow de GitHub Actions para releases automáticas en push a `main`
  - Detecta la versión desde el `CHANGELOG.md` automáticamente
  - Actualiza el badge y links de versión en `README.md`
  - Compila binarios para Windows, Linux, macOS Intel y Apple Silicon
  - Sube binarios y ZIP del código fuente a GitHub Releases
  - Incluye las notas del CHANGELOG en el body de la release

### Archivos añadidos

- `cmd/version.go` — nuevo comando `version`
- `.github/workflows/release.yml` — workflow de release automática

### Archivos modificados

- `main.go` — registro de `VersionCmd`
- `docs/commands.md` — documentación del comando `version`

---

## [1.9.0] - 2026-06-11 — "Oe, que buenas camaras csm"

### Añadido

- Nuevo comando `wachiman monitor` — supervisor activo de contenedores
  - Detecta cambios de estado en tiempo real con intervalo configurable (`--interval`)
  - Modo activo: reinicia automáticamente contenedores caídos
  - Modo observación (`--no-restart`): solo notifica sin actuar
  - Filtro por contenedor (`--only nginx,db`) para monitoreo selectivo
  - Anti-spam: compara estado simplificado (`running`/`stopped`) en vez del string completo del status
- Nuevo comando `wachiman network` con subcomandos para gestión de redes
  - `wachiman network ls` — lista todas las redes con driver y scope
  - `wachiman network inspect [nombre]` — detalles de una red: subnet, gateway, contenedores con IPs
  - `wachiman network connect [red] [contenedor]` — conecta un contenedor a una red
  - `wachiman network disconnect [red] [contenedor]` — desconecta un contenedor de una red
- Progreso visual en `wachiman backup` con spinner, velocidad de transferencia y tamaño final
  - Usa `io.MultiWriter` para escribir al archivo y actualizar la barra simultáneamente
  - Sin tamaño conocido de antemano — spinner en modo indeterminado

### Archivos modificados

- `cmd/monitor.go` — nuevo comando de supervisión activa
- `cmd/network.go` — nuevo comando de gestión de redes con 4 subcomandos
- `cmd/backup.go` — progreso visual con `schollz/progressbar`
- `docker/client.go` — nuevas funciones `ListNetworks`, `InspectNetwork`, `NetworkConnect`, `NetworkDisconnect`
- `docs/commands.md` — documentación de `monitor` y `network`
- `main.go` — registro de `MonitorCmd` y `NetworkCmd`

---

## [1.5.0] - 2026-06-10 — "Llaves y linterna, supongo que para guardar las cosas en la bodega"

### Añadido
- Nuevo comando `wachiman backup [nombre_contenedor]` para crear respaldos empaquetados en `.tar`.
  - **Pausado inteligente (`Pause`/`Unpause`):** Pausa automáticamente el contenedor antes de copiar para prevenir la corrupción de datos y asegura su reactivación mediante `defer`, incluso si el proceso falla.
  - **Filtro de rutas redundantes:** Algoritmo de optimización que analiza los volúmenes montados y descarta subcarpetas o subarchivos si su directorio raíz ya va a ser respaldado (ej. no duplica `/var/www/html/wp-content` si ya está respaldando `/var/www/html`).
  - **Modo Hot Backup (`--no-pause`):** Flag para forzar el respaldo en caliente sin detener el contenedor (advirtiendo sobre posibles inconsistencias).
- Nuevo comando `wachiman export-compose [nombre_contenedor]` para realizar ingeniería inversa y generar dinámicamente un manifiesto `docker-compose.yml` funcional mapeando puertos, variables de entorno limpias y volúmenes locales.
- Nuevo comando `wachiman shell [nombre_contenedor]` para abrir una terminal interactiva dentro del contenedor de forma rápida, intentando usar `bash` y cayendo en `sh` si no está disponible.
- Nuevo comando `wachiman audit [nombre_contenedor]` para realizar un análisis rápido de seguridad y optimización del contenedor (revisión de puertos expuestos, variables de entorno sensibles y usuarios ejecutores).

### Cambios / Mejoras
- Corrección y actualización del cliente en `docker/client.go` para aislar la estructura interna de `Inspect`, protegiendo el CLI contra cambios drásticos de firmas y tipos en las actualizaciones del SDK oficial de Docker.

### Archivos modificados
- `cmd/backup.go` — nuevo comando de copias de seguridad con desduplicación de rutas.
- `cmd/export_compose.go` — nuevo comando para ingeniería inversa a docker-compose.
- `cmd/shell.go` — nuevo comando para acceso interactivo por terminal (TTY).
- `cmd/audit.go` — nuevo comando de auditoría y buenas prácticas.
- `docker/client.go` — refactorización del método `Inspect` y abstracción de la estructura de volúmenes.
- `main.go` — registro de los nuevos comandos en el CLI raíz.

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

---
## [0.2.0] - 2026-06-08 — "Seño, pan con palta y su quinua con manzana"

### Añadido

- `wachiman watch` ahora muestra sparklines de tendencia con caracteres Unicode `▁▂▃▄▅▆▇█`
  - Historial de los últimos 10 ticks por contenedor en memoria
  - Color según el último valor: verde < 50%, amarillo 50–80%, rojo > 80%
- Header dinámico en `wachiman watch` con conteo de contenedores corriendo vs parados
- Fix de pantalla en Windows — `cls` en vez de escape codes ANSI para limpiar correctamente

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

