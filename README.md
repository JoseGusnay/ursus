# Ursus 🐻
> El cerebro externo para tus agentes de IA.

**Ursus** es un sistema de memoria persistente de nivel profesional diseñado para agentes de IA. Permite que tus asistentes (Claude, Cursor, Windsurf) recuerden decisiones, aprendizajes y contextos a través de múltiples sesiones de trabajo.

![Ursus Banner](https://img.shields.io/badge/Status-100%25_Parity_with_Engram-blue?style=for-the-badge&logo=go)

---

## 🛠️ Instalación (¿Cómo lo uso?)

Para un usuario final, Ursus se instala de forma global en el sistema para que esté disponible en cualquier terminal o proyecto.

### 1. Requisitos
- **Go 1.22** o superior.

### 2. Instalación Global
Ejecuta el siguiente comando para instalar el binario `ursus` en tu carpeta `GOBIN`:
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

## 🤖 Integración con Agentes (MCP)

Ursus es un servidor **MCP (Model Context Protocol)**. Para que un agente como **Claude Desktop** lo use, añade esto a tu archivo de configuración (`claude_desktop_config.json`):

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
