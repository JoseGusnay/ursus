# 🚀 Guía de Instalación (Usuario Final)

Ursus está diseñado para ser instalado en segundos, incluso si no eres desarrollador.

## Instalación en Windows (Zero-Go)

1. **Descarga**: Obtén el archivo `ursus_windows_amd64.exe` desde [Releases](https://github.com/JoseGusnay/ursus/releases).
2. **Ubicación**: Crea una carpeta (ej. `C:\Herramientas\ursus`) y mueve el archivo allí.
3. **Cambio de nombre**: Renombra el archivo a `ursus.exe`.
4. **Instalación Automática**: Abre una terminal (CMD o PowerShell) en esa carpeta y ejecuta:
   ```powershell
   .\ursus.exe setup path
   ```
   *Esto configurará las variables de entorno de Windows por ti.*
5. **Reinicio**: Cierra la terminal y abre una nueva. Ahora el comando `ursus` funcionará desde cualquier carpeta.

## Solución de Problemas

### Aviso de SmartScreen
Al ser una aplicación de código abierto, Windows puede mostrar un aviso de "PC Protegido".
- Haz clic en **"Más información"**.
- Haz clic en **"Ejecutar de todas formas"**.

### Comando no reconocido
Si después de hacer el `setup path` y reiniciar la terminal el comando sigue sin funcionar:
1. Asegúrate de cerrar **todas** las ventanas de la terminal.
2. Verifica manualmente que la ruta de tu carpeta esté en las "Variables de Entorno" de tu cuenta de usuario.

---
[Volver al README](../README.md)
