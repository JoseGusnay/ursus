# 💻 Referencia de Comandos CLI

Ursus ofrece una interfaz de línea de comandos potente para gestionar tu memoria manualmente.

## Gestión de Memorias

### `ursus add "[contenido]"`
Añade una memoria rápida.
- `--topic, -t`: Categoriza la memoria (ej. `--topic auth`).
- `--metadata, -m`: Añade datos extra en formato JSON o texto.

### `ursus search "[consulta]"`
Búsqueda ultrarrápida a texto completo. Encuentra memorias por palabras clave o temas.

### `ursus list`
Muestra las memorias más recientes.
- `--limit, -l`: Limita el número de resultados (ej. `ursus list -l 5`).

### `ursus delete [id]`
Borra una memoria por su ID (Borrado lógico/Soft Delete).

### `ursus update [id] "[nuevo contenido]"`
Actualiza el contenido de una memoria existente.

## Análisis y Visualización

### `ursus stats`
Muestra un reporte de tu actividad: cuántas memorias tienes, temas más usados y actividad en la última semana.

### `ursus tui`
Abre una interfaz gráfica en la terminal (TUI) para navegar por tus memorias con estilo.

### `ursus timeline`
Muestra tus memorias agrupadas por fecha en orden cronológico.

## Sesiones de Trabajo

### `ursus session start "[título]"`
Inicia una sesión para agrupar aprendizajes específicos de una tarea.

### `ursus session end`
Cierra la sesión activa.

### `ursus review`
Muestra un resumen de la sesión activa o de una sesión específica por su ID.

## Configuración

### `ursus setup path`
Añade Ursus al PATH de Windows.

### `ursus setup claude`
Configura Ursus automáticamente en Claude Desktop.

---
[Volver al README](../README.md)
