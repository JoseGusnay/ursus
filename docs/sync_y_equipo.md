# 🔄 Sincronización y Trabajo en Equipo

Ursus permite que la memoria técnica viaje con el código fuente a través de Git, ideal para equipos que quieren compartir contexto.

## El Comando `sync`

El flujo de trabajo recomendado es:

1. **Durante el trabajo**: Tu agente guarda memorias en tu base de datos global personal.
2. **Al terminar**: Ejecutas `ursus sync`.
    - Esto extrae los aprendizajes relevantes del proyecto actual.
    - Los empaqueta en pequeños trozos comprimidos en la carpeta `.ursus/`.
3. **Commit & Push**: Sube la carpeta `.ursus/` a tu repositorio Git.

## Cómo lo usa tu equipo

Cuando otro desarrollador clona el repo:
1. Ejecuta `ursus sync`.
2. Ursus detecta los archivos en `.ursus/` y los importa a su base de datos local.
3. ¡Su IA ahora tiene todo el contexto histórico del proyecto sin haber estado presente en las reuniones originales!

## Detalles Técnicos
- **Formato**: Usamos Chunks de JSONL comprimidos con Gzip para no ensuciar los diffs de Git.
- **Eficiencia**: Solo se sincronizan memorias marcadas con el "scope" del proyecto actual.

---
[Volver al README](../README.md)
