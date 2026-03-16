# 🤖 Guía Maestra de MCP (Model Context Protocol)

Ursus brilla cuando se conecta a tus agentes de IA. Utiliza el estándar industrial **MCP** para actuar como el cerebro externo de tus asistentes.

## Herramientas Disponibles (Tools)

Cuando conectas Ursus a un agente (Cursor, Claude, etc.), este gana acceso a las siguientes herramientas:

| Herramienta | Uso |
| :--- | :--- |
| `add_memory` | El agente guarda un aprendizaje nuevo de la conversación actual. |
| `search_memory` | El agente busca en tus memorias pasadas para responder con contexto. |
| `passive_capture` | El agente extrae aprendizajes automáticamente de su propia respuesta. |
| `mem_stats` | El agente te da un reporte de cuánto ha aprendido de ti. |
| `session_start/end` | El agente agrupa aprendizajes en una sesión de trabajo específica. |
| `summarize_session` | El agente genera un resumen ejecutivo de lo logrado en una sesión. |

## Configuración por Agente

### Cursor / Windsurf
1. Ve a **Settings -> Features -> MCP**.
2. Añade un nuevo servidor:
   - **Name**: `Ursus`
   - **Type**: `command`
   - **Command**: `ursus mcp`

### Claude Desktop
Ejecuta el comando de configuración automática:
```bash
ursus setup claude
```

## Consejos de Uso
- **Pídele a tu IA que busque**: "Busca en Ursus lo que decidimos sobre la arquitectura la semana pasada".
- **Deja que aprenda sola**: Ursus está configurado para que, si la IA ve una sección de "Aprendizajes" en su respuesta, la guarde automáticamente.

---
[Volver al README](../README.md)
