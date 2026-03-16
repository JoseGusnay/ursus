# 🛡️ Privacidad y Seguridad (Local-First)

Tu conocimiento es tuyo. Ursus está diseñado bajo la filosofía **"Privacidad por Defecto"**.

## Motor de Redacción Automática

Ursus incluye un procesador de lenguaje que limpia tus datos **antes** de guardarlos en la base de datos local. Esto evita que información sensible "contamine" tu memoria a largo plazo.

### Qué se oculta automáticamente:
- **Correos Electrónicos**: `usuario@ejemplo.com` -> `[EMAIL_REDACTED]`
- **Secretos y Tokens**:
    - Keys de OpenAI (`sk-...`)
    - Tokens de autenticación (`secret_...`, `token-...`)
    - Claves generales (`key-...`)
- **Contenido Manual**: Cualquier texto envuelto en etiquetas `<private>...</private>` será reemplazado por `[REDACTED]`.

## Seguridad Local-First

1. **Sin Nube**: Ursus no envía tus datos a ningún servidor externo. Todo vive en tu máquina bajo tu control.
2. **Procesamiento en el Binario**: La redacción de privacidad ocurre en el binario de Go en tu hardware, mucho antes de que se escriba en el disco o se envíe a un agente de IA.
3. **SQLite con FTS5**: Tus datos se guardan en un archivo `.db` estándar que puedes auditar o borrar en cualquier momento en `~/.ursus/ursus.db`.

---
[Volver al README](../README.md)
