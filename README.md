# Wachiman CLI

> **Wachiman** (del inglés watchman) • *Sustantivo (Perú).* Vigilante, guardián.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Engine%20%2F%20Desktop-2496ED?style=flat-square&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg?style=flat-square)](LICENSE)
[![Latest Release](https://img.shields.io/badge/Release-v1.9.5-green?style=flat-square)](https://github.com/x-name15/wachiman-cli/releases)

![wachiman](images/wachiman.png)

---

## Requisitos

- Go 1.21+
- Docker Desktop o Docker Engine corriendo localmente

---

## Instalación

**Build local:**
```bash
go build -o wachiman .
```

**Instalación global** (añade `wachiman` a tu PATH via `$GOPATH/bin`):
```bash
go install .
```

**Windows:**
```bash
go build -o wachiman.exe .
```
---

## Documentación

- Guía de comandos: [docs/commands.md](docs/commands.md)

---

## Construido con

- [cobra](https://github.com/spf13/cobra) — framework para CLIs en Go
- [moby/moby](https://github.com/moby/moby) — SDK oficial de Docker para Go
- [fatih/color](https://github.com/fatih/color) — colores en terminal

---

## Licencia

GNU GENERAL PUBLIC LICENSE Version 3

---

## Changelog

Se añadió un registro de cambios (changelog) a este `README.md`. Consulta la última versión y notas de la release en:
- [v1.9.5 - CHANGE LOG](CHANGELOG.md)
- [v1.9.5 - Release notes](https://github.com/x-name15/wachiman-cli/releases/tag/1.9.5)
