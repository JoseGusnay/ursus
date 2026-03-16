# 🐻 Ursus: Tu Memoria Inteligente para Agentes de IA

**Ursus** es un sistema de memoria persistente de nivel profesional diseñado para agentes de IA (Cursor, Windsurf, Claude Desktop, etc.). Permite que tus asistentes recuerden decisiones, arquitecturas y lecciones aprendidas a través de múltiples sesiones y proyectos.

Inspirado en Engram, pero con una arquitectura robusta en Go, optimizada para velocidad y privacidad.

---

## 🛠️ Instalación (¿Cómo lo uso?)

Para un "Usuario Final", Ursus se distribuye como un **binario único y autocontenido**. No necesitas instalar Go ni configurar variables del sistema a mano.

### Opción A: Descarga Directa (Recomendado)
1. Ve a [GitHub Releases](https://github.com/JoseGusnay/ursus/releases).
2. Descarga `ursus_windows_amd64.exe` (o la versión para Mac/Linux).
3. **¡Instalación Automática!**: Abre una terminal en la carpeta donde lo descargaste y ejecuta:
   ```powershell
   .\ursus_windows_amd64.exe setup path
   ```
   *Esto añade Ursus a tu PATH de Windows automáticamente. Cierra y abre la terminal para activarlo.*
4. **Cero Complicaciones**: Ahora puedes renombrar el archivo a `ursus.exe` y usarlo desde cualquier lugar simplemente escribiendo `ursus`.

> [!NOTE]
> **Aviso de SmartScreen**: Al abrirlo por primera vez, haz clic en "Más información" -> "Ejecutar de todas formas". Es normal por ser código abierto recién compilado.

### Opción B: Para Desarrolladores (Go Install)
```bash
go install github.com/JoseGusnay/ursus/cmd/ursus@latest
```

---

## 🤖 Configuración de Agentes (Agent Setup)

### 1. Claude Desktop
Configúralo automáticamente con un comando:
```bash
ursus setup claude
```

### 2. Cursor / Windsurf / VS Code
1. Ve a **Settings -> Features -> MCP**.
2. Añade un nuevo servidor:
   - **Name**: `ursus`
   - **Type**: `command`
   - **Command**: `ursus mcp`

---

## 🏗️ Modelo de Memoria Híbrido

Ursus combina lo mejor de dos mundos:
1.  **Memoria Global**: Guardada en `~/.ursus/ursus.db`. El agente te conoce a ti y tus preferencias en cualquier PC o proyecto.
2.  **Memoria de Proyecto**: Sincroniza conocimientos específicos del repositorio mediante archivos `.jsonl.gz` Git-friendly usando `ursus sync`.

---

## 🎮 Comandos de la CLI

| Comando | Descripción |
| :--- | :--- |
| `ursus stats` | Ver tu actividad, sesiones y temas recurrentes. |
| `ursus tui` | Abre la interfaz visual azul interactiva. |
| `ursus search` | Búsqueda semántica ultrarrápida en tu cerebro. |
| `ursus add` | Guarda una memoria manual con tópicos. |
| `ursus review` | Resume lo ocurrido en la sesión actual. |
| `ursus setup` | Automatiza PATH y configuración de agentes. |

---

## ✨ Diferenciadores Clave
- **Privacidad**: Filtra automáticamente Tokens, API Keys y datos sensibles.
- **Higiene**: Deduplicación dinámica (no guarda basura repetida).
- **Velocidad**: Basado en SQLite con FTS5 para búsquedas instantáneas.
- **Portabilidad**: Un solo archivo para todo.

---
Developed with ❤️ by **Jose Gusnay & Antigravity**
