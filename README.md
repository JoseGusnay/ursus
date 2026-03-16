# Ursus 🐻
> El cerebro externo para tus agentes de IA.

**Ursus** es un sistema de memoria persistente de nivel profesional diseñado para agentes de IA. Permite que tus asistentes (Claude, Cursor, Windsurf) recuerden decisiones, aprendizajes y contextos a través de múltiples sesiones de trabajo.

![Ursus Banner](https://img.shields.io/badge/Status-100%25_Parity_with_Engram-blue?style=for-the-badge&logo=go)

---

## 🛠️ Instalación (¿Cómo lo uso?)

Ursus se distribuye como un **binario único y autocontenido**. Esto significa que **NO necesitas tener instalado Go** para usarlo.

### Opción A: Descargar el Binario (Recomendado)
1. Ve a la sección de [Releases de GitHub](https://github.com/JoseGusnay/ursus/releases).
2. Descarga la versión correspondiente a tu sistema operativo (Windows, Linux o macOS).
3. Descomprime el archivo y coloca el ejecutivo en una carpeta que esté en tu PATH (ej. `C:\Users\TuUsuario\bin\`).

### Opción B: Instalación para Desarrolladores (Requiere Go)
Si eres desarrollador y tienes Go instalado, puedes instalarlo directamente:
```bash
go install github.com/JoseGusnay/ursus/cmd/ursus@latest
```

### 3. Verificar instalación
```bash
ursus stats
```

---

## 🏗️ Modelo de Memoria Híbrido

Ursus resuelve el problema de "¿dónde se guardan los datos?" con un enfoque inteligente:

1.  **Memoria Global (Personal)**: Por defecto, Ursus guarda todo en una base de datos maestra en `~/.ursus/ursus.db`. Esto permite que el agente te conozca a ti y a tus preferencias generales, sin importar en qué proyecto estés.
2.  **Memoria de Proyecto (Local)**: Ursus permite "sincronizar" memorias con un proyecto específico. Al usar el comando `ursus sync`, el sistema empaqueta el conocimiento relevante en una carpeta `.ursus/` dentro de tu repositorio. Esto permite que otros desarrolladores del mismo proyecto compartan el contexto.

---

## 🤖 Configuración de Agentes (Agent Setup)

Para que Ursus sea útil, debes conectarlo a tus herramientas favoritas. Sigue estos pasos según tu agente:

### 1. Claude Desktop
Añade lo siguiente a tu archivo `claude_desktop_config.json` (ubicado en `%APPDATA%\Claude\config\` en Windows):

```json
{
  "mcpServers": {
    "ursus": {
      "command": "ursus",
      "args": ["mcp"]
    }
  }
}
```

### 2. Cursor / Windsurf
1. Ve a **Settings** -> **Features** -> **MCP**.
2. Añade un nuevo servidor:
   - **Name**: `ursus`
   - **Type**: `command`
   - **Command**: `ursus mcp`

### 3. VS Code (Copilot / Claude Code)
Si usas extensiones que consumen MCP:
1. Asegúrate de que `ursus` esté en tu `$PATH`.
2. Registra el comando `ursus mcp` en la sección de servidores MCP de la extensión.

---

## 🎮 Formas de Uso

### 1. Interfaz TUI (Exploración Visual)
Ideal para revisar qué ha aprendido el agente de forma rápida y estética.
```bash
ursus tui
```

### 2. Gestión CLI (Control Total)
- **Añadir aprendizaje**: `ursus add "El backend usa Clean Architecture" --topic "arch"`
- **Ver estadísticas**: `ursus stats`
- **Sincronizar con el repo**: `ursus sync`

### 3. Captura Pasiva
El agente puede guardar memorias automáticamente si en su respuesta incluye etiquetas como:
```markdown
### Aprendizajes
- La configuración de CORS debe ir antes de las rutas.
```

---

## 🧩 Diferenciadores Clave
- **Privacidad**: Filtra automáticamente contenido sensible marcado con `<private>`.
- **Higiene**: Deduplicación automática para no llenar la DB de basura.
- **Arquitectura**: Diseñado bajo **Clean Architecture**, lo que garantiza estabilidad y facilidad de expansión.

---
Developed with ❤️ by **Jose Gusnay & Antigravity**
