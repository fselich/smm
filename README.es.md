**[English](README.md)** | **Español**

# SMM - Secret Manager Manager

**SMM** es una herramienta de interfaz de terminal (TUI) que permite visualizar, editar y gestionar secretos de Google Cloud Platform de manera eficiente y segura.

## Características

- **Navegación intuitiva** con interfaz de terminal moderna
- **Edición de secretos** con tu editor favorito
- **Gestión de versiones** - visualiza, restaura y crea nuevas versiones
- **Búsqueda avanzada** - busca por nombre o contenido
- **Copia al portapapeles** con un solo comando
- **Syntax highlighting** para múltiples formatos (JSON, Bash, INI, PHP)
- **Multi-proyecto** - cambia fácilmente entre proyectos de GCP

## Instalación

### Homebrew (Recomendado)

```bash
brew tap fselich/tap
brew install smm
```

### Descargar Binario

Descarga la última versión para tu plataforma desde [GitHub Releases](https://github.com/fselich/smm/releases):

#### Linux x64
```bash
wget https://github.com/fselich/smm/releases/latest/download/smm-linux-amd64.tar.gz
tar -xzf smm-linux-amd64.tar.gz
chmod +x smm
sudo mv smm /usr/local/bin/  # Opcional: instalar en el sistema
```

#### Linux ARM64
```bash
wget https://github.com/fselich/smm/releases/latest/download/smm-linux-arm64.tar.gz
tar -xzf smm-linux-arm64.tar.gz
chmod +x smm
sudo mv smm /usr/local/bin/  # Opcional: instalar en el sistema
```

### Desde el código fuente

#### Requisitos previos
- [Go](https://go.dev/doc/install) (versión 1.24.4 o superior)
- [Git](https://git-scm.com/book/en/v2/Getting-Started)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) (gcloud)

#### Clonar el repositorio y compilar
```bash
git clone git@github.com:fselich/smm.git
cd smm
go build -o smm cmd/main.go
```

## Uso

### Uso básico
```bash
./smm                    # Usar proyecto por defecto
./smm -p PROJECT_ID      # Especificar proyecto de GCP
./smm -v                 # Mostrar información de la versión
```
## Controles de Teclado

### Navegación
| Tecla       | Acción                                                     |
| ----------- | ---------------------------------------------------------- |
| `↑` `↓`     | Navegar por la lista / Scroll en el detalle del secreto   |
| `Tab`       | Cambiar foco entre lista y detalle                        |
| `Shift + ←→`| Redimensionar la vista de la lista                        |

### Búsqueda y Filtrado
| Tecla       | Acción                                                     |
| ----------- | ---------------------------------------------------------- |
| `/`         | Filtrar por nombre de secreto                              |
| `Ctrl+F`    | Buscar en el contenido de todos los secretos              |

### Gestión de Secretos
| Tecla       | Acción                                                     |
| ----------- | ---------------------------------------------------------- |
| `i`         | Mostrar información del secreto (metadatos, fecha de creación, etiquetas) |
| `c`         | Copiar secreto al portapapeles                             |
| `n`         | Crear nueva versión del secreto                            |
| `v`         | Mostrar/ocultar versiones del secreto                      |
| `r`         | Restaurar versión seleccionada                             |

### Sistema
| Tecla       | Acción                                                     |
| ----------- | ---------------------------------------------------------- |
| `p`         | Cambiar proyecto de GCP                                    |
| `Esc`       | Refrescar / Cancelar operación                             |
| `Ctrl+C`    | Salir del programa                                         |

### Opciones de línea de comandos

| Opción            | Descripción                                    |
| ----------------- | ---------------------------------------------- |
| `-p PROJECT_ID`   | Cargar secretos del proyecto especificado     |
| `-v`              | Mostrar información de la versión             |


## Autenticación

SMM utiliza la autenticación existente de `gcloud`. Asegúrate de estar autenticado antes de usar la herramienta.

### Verificar autenticación

```bash
gcloud config list
```

El resultado debe contener tu cuenta y proyecto:

```bash
[core]
account = tu@email.com
project = tu-proyecto-gcp
```

### Verificar permisos

```bash
gcloud secrets list --project=tu-proyecto-gcp
```

Si tienes los permisos necesarios, verás la lista de secretos del proyecto.

### Permisos requeridos

Tu cuenta necesita los siguientes roles de IAM:
- `roles/secretmanager.viewer` - Para listar y leer secretos
- `roles/secretmanager.secretVersionManager` - Para crear nuevas versiones

### Autenticación con gcloud
Si no estás autenticado, puedes hacerlo con:
```bash
gcloud auth login
```

## Syntax Highlighting

SMM detecta automáticamente el formato del contenido y aplica coloreado de sintaxis para:

- **Bash/Env** - Variables de entorno
- **JSON** - Datos estructurados  
- **INI** - Archivos de configuración
- **PHP** - Código PHP

## Configuración

La aplicación almacena su configuración en `~/.config/smm/config.yaml`. **El archivo de configuración se crea automáticamente si no existe** cuando ejecutas SMM por primera vez.

Estructura de configuración:

```yaml
projects:                               # Lista de proyectos GCP
  - id: "mi-proyecto-gcp-1"
    type: "gcp"
  - id: "mi-proyecto-gcp-2"
    type: "gcp"  
selected: "mi-proyecto-gcp-1"            # Proyecto actualmente seleccionado
logPath: "/ruta/al/archivo/log"          # Ruta del archivo de log (opcional)
```

**Notas:**
- Los proyectos se añaden automáticamente cuando cambias a ellos usando la tecla `p`
- El campo `selected` recuerda tu último proyecto usado
- `logPath` es opcional - déjalo vacío para deshabilitar el logging

## Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/nueva-caracteristica`)
3. Commit tus cambios (`git commit -am 'Añadir nueva característica'`)
4. Push a la rama (`git push origin feature/nueva-caracteristica`)  
5. Abre un Pull Request

## Licencia

Este proyecto está bajo la licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.
