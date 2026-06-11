# Comandos

Referencia completa de comandos disponibles en **Wachiman**.

---

## `ps` — Lista contenedores

Muestra los contenedores Docker detectados por el sistema.

```bash
wachiman ps
```

Solo contenedores activos:

```bash
wachiman ps --running
```

Solo contenedores detenidos:

```bash
wachiman ps --stopped
```

Salida en JSON:

```bash
wachiman ps --output json
```

Opciones:

| Flag           | Descripción                               |
| -------------- | ----------------------------------------- |
| `--running`    | Mostrar únicamente contenedores activos   |
| `--stopped`    | Mostrar únicamente contenedores detenidos |
| `-o, --output` | Formato de salida (`table`, `json`)       |

---

## `stats` — Uso de CPU y memoria

Muestra estadísticas de recursos de los contenedores activos.

```bash
wachiman stats
```

Salida en JSON:

```bash
wachiman stats --output json
```

Opciones:

| Flag           | Descripción                         |
| -------------- | ----------------------------------- |
| `-o, --output` | Formato de salida (`table`, `json`) |

---

## `watch` — Monitor en tiempo real

Muestra un dashboard actualizado periódicamente con información de los contenedores.

```bash
wachiman watch
```

Intervalo personalizado:

```bash
wachiman watch --interval 5
wachiman watch -i 10
```

Opciones:

| Flag             | Descripción                       |
| ---------------- | --------------------------------- |
| `-i, --interval` | Intervalo de refresco en segundos |

---

## `logs` — Logs de un contenedor

Muestra los registros de un contenedor.

```bash
wachiman logs wordpress_app
```

Mostrar más líneas:

```bash
wachiman logs wordpress_app --tail 100
```

Forma corta:

```bash
wachiman logs wordpress_app -t 100
```

Opciones:

| Flag         | Descripción                |
| ------------ | -------------------------- |
| `-t, --tail` | Número de líneas a mostrar |

---

## `inspect` — Inspecciona un contenedor

Muestra información detallada de un contenedor.

```bash
wachiman inspect wordpress_app
```

Incluye información como:

* Estado
* Imagen utilizada
* Dirección IP
* Puertos expuestos
* Volúmenes montados
* Variables de entorno

---

## `top` — Procesos dentro del contenedor

Lista los procesos que se están ejecutando dentro del contenedor.

```bash
wachiman top wordpress_app
```

---

## `shell` — Terminal interactiva

Abre una sesión dentro del contenedor.

```bash
wachiman shell wordpress_app
```

Seleccionar shell específica:

```bash
wachiman shell wordpress_app --shell /bin/bash
```

Forma corta:

```bash
wachiman shell wordpress_app -s /bin/bash
```

Opciones:

| Flag          | Descripción      |
| ------------- | ---------------- |
| `-s, --shell` | Shell a utilizar |

---

## `start` — Arranca un contenedor

Inicia un contenedor detenido.

```bash
wachiman start wordpress_app
```

---

## `stop` — Detiene un contenedor

Detiene un contenedor en ejecución.

```bash
wachiman stop wordpress_app
```

---

## `restart` — Reinicia un contenedor

Reinicia un contenedor existente.

```bash
wachiman restart wordpress_app
```

---

## `backup` — Respaldo de datos

Genera un archivo `.tar` con la información asociada al contenedor.

```bash
wachiman backup wordpress_app
```

Guardar en un directorio específico:

```bash
wachiman backup wordpress_app --output ./backups
```

Forma corta:

```bash
wachiman backup wordpress_app -o ./backups
```

Realizar respaldo sin pausar el contenedor:

```bash
wachiman backup wordpress_app --no-pause
```

Opciones:

| Flag           | Descripción                                                 |
| -------------- | ----------------------------------------------------------- |
| `--no-pause`   | Realiza el respaldo sin detener temporalmente el contenedor |
| `-o, --output` | Directorio donde guardar el archivo `.tar`                  |

---

## `network` — Gestiona las redes de Docker

### Listar redes

```bash
wachiman network ls
```

### Inspeccionar una red

Muestra detalles de una red — subnet, gateway, y contenedores conectados con sus IPs.

```bash
wachiman network inspect bridge
wachiman network inspect wordpress_default
```

### Desconectar un contenedor de una red

```bash
wachiman network disconnect [red] [contenedor]
```

Ejemplo:

```bash
wachiman network disconnect wordpress_default wordpress_app
```

### Conectar un contenedor a una red

```bash
wachiman network connect [red] [contenedor]
```

Ejemplo:

```bash
wachiman network connect bridge wordpress_app
```

Subcomandos disponibles:

| Subcomando   | Descripción                                    |
| ------------ | ---------------------------------------------- |
| `ls`         | Lista todas las redes                          |
| `inspect`    | Muestra detalles de una red específica         |
| `connect`    | Conecta un contenedor a una red                |
| `disconnect` | Desconecta un contenedor de una red            |

> ⚠️ Desconectar un contenedor de su red puede interrumpir su comunicación con otros servicios.

---
## `audit` — Auditoría de seguridad

Analiza configuraciones potencialmente inseguras de un contenedor.

```bash
wachiman audit wordpress_app
```

La auditoría puede detectar:

* Contenedores privilegiados
* Puertos expuestos
* Configuraciones peligrosas
* Variables sensibles
* Permisos excesivos

---

## `export-compose` — Exportar a Docker Compose

Genera un archivo `docker-compose.yml` basado en la configuración actual del contenedor.

```bash
wachiman export-compose wordpress_app
```

---

## `prune` — Limpieza de recursos

Elimina recursos Docker sin uso.

Incluye:

* Contenedores detenidos
* Imágenes huérfanas
* Volúmenes sin uso

```bash
wachiman prune
```

Sin solicitar confirmación:

```bash
wachiman prune --force
```

Forma corta:

```bash
wachiman prune -f
```

Opciones:

| Flag          | Descripción                                |
| ------------- | ------------------------------------------ |
| `-f, --force` | Ejecuta la limpieza sin pedir confirmación |

> ⚠️ Esta operación es irreversible.

---

## `config` — Configuración global

Gestiona la configuración de Wachiman.

### Ver configuración actual

```bash
wachiman config get
```

### Modificar una configuración

```bash
wachiman config set <clave> <valor>
```

Ejemplo:

```bash
wachiman config set refresh_interval 5
```

### Restaurar valores por defecto

```bash
wachiman config reset
```

Subcomandos disponibles:

| Comando | Descripción                           |
| ------- | ------------------------------------- |
| `get`   | Muestra la configuración actual       |
| `set`   | Modifica una configuración            |
| `reset` | Restaura la configuración por defecto |

---

## `monitor` — Monitorea y reinicia contenedores caídos

Vigila el estado de los contenedores en tiempo real. Si detecta que uno se cayó, lo reinicia automáticamente.

```bash
wachiman monitor
```

Solo observar sin reiniciar:

```bash
wachiman monitor --no-restart
```

Intervalo personalizado:

```bash
wachiman monitor --interval 10
wachiman monitor -i 10
```

Monitorear contenedores específicos:

```bash
wachiman monitor --only wordpress_app,wordpress_db
```

Opciones:

| Flag            | Descripción                                        |
| --------------- | -------------------------------------------------- |
| `-i, --interval`| Intervalo de chequeo en segundos (por defecto: 5)  |
| `--no-restart`  | Solo observar, no reiniciar contenedores caídos    |
| `--only`        | Contenedores a monitorear (separados por coma)     |

> ⚠️ El modo activo reinicia automáticamente cualquier contenedor que se detenga. Usa `--no-restart` si solo quieres observar.

---
## Ayuda

Ayuda general:

```bash
wachiman --help
```

Ayuda de un comando específico:

```bash
wachiman logs --help
wachiman stats --help
wachiman watch --help
```
