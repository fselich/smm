**[🇺🇸 English](README.md)** | **🇪🇸 Español**

# SMM - Secret Manager Manager

**SMM** es una herramienta de interfaz de terminal (TUI) que permite visualizar, editar y gestionar secretos de Google Cloud Platform de manera eficiente y segura.

## ✨ Características

- 🔍 **Navegación intuitiva** con interfaz de terminal moderna
- 📝 **Edición de secretos** con tu editor favorito
- 🔄 **Gestión de versiones** - visualiza, restaura y crea nuevas versiones
- 🔎 **Búsqueda avanzada** - busca por nombre o contenido
- 📋 **Copia al portapapeles** con un solo comando
- 🎨 **Syntax highlighting** para múltiples formatos (JSON, Bash, INI, PHP)
- 🚀 **Multi-proyecto** - cambia fácilmente entre proyectos de GCP

## 📦 Instalación

### Desde el código fuente

#### Requisitos previos
- [Go](https://go.dev/doc/install) (versión 1.24.4 o superior)
- [Git](https://git-scm.com/book/en/v2/Getting-Started)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) (gcloud)
- libx11-dev para soporte del portapapeles en Linux (`sudo apt-get install libx11-dev` o similar)

#### Clonar el repositorio y compilar
```bash
git clone git@github.com:fselich/smm.git
cd smm
go build -o smm cmd/main.go
```

## 🚀 Uso

### Uso básico
```bash
./smm                    # Usar proyecto por defecto
./smm -p PROJECT_ID      # Especificar proyecto de GCP
```
## ⌨️ Controles de Teclado

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


## 🔐 Autenticación

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

## 🎨 Syntax Highlighting

SMM detecta automáticamente el formato del contenido y aplica coloreado de sintaxis para:

- 🌱 **Bash/Env** - Variables de entorno
- 📄 **JSON** - Datos estructurados  
- ⚙️ **INI** - Archivos de configuración
- 🐘 **PHP** - Código PHP

## 📁 Configuración

La aplicación almacena su configuración en `~/.config/smm/config.yaml`:

```yaml
projectIds: ["proyecto-1", "proyecto-2"]  # Proyectos disponibles
selected: "proyecto-1"                    # Proyecto seleccionado
logPath: "/path/to/log/file"             # Archivo de log (opcional)
```

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/nueva-caracteristica`)
3. Commit tus cambios (`git commit -am 'Añadir nueva característica'`)
4. Push a la rama (`git push origin feature/nueva-caracteristica`)  
5. Abre un Pull Request

## 📝 Licencia

Este proyecto está bajo la licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.
