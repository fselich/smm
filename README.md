**🇺🇸 English** | **[🇪🇸 Español](README.es.md)**

# SMM - Secret Manager Manager

**SMM** is a Terminal User Interface (TUI) tool that allows you to efficiently and securely view, edit, and manage Google Cloud Platform secrets.

## ✨ Features

- 🔍 **Intuitive navigation** with modern terminal interface
- 📝 **Secret editing** with your favorite editor
- 🔄 **Version management** - view, restore, and create new versions
- 🔎 **Advanced search** - search by name or content
- 📋 **Copy to clipboard** with a single command
- 🎨 **Syntax highlighting** for multiple formats (JSON, Bash, INI, PHP)
- 🚀 **Multi-project** - easily switch between GCP projects

## 📦 Installation

### From source

#### Prerequisites
- [Go](https://go.dev/doc/install) (version 1.24.4 or higher)
- [Git](https://git-scm.com/book/en/v2/Getting-Started)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) (gcloud)

#### Clone repository and build
```bash
git clone git@github.com:fselich/smm.git
cd smm
go build -o smm cmd/main.go
```

## 🚀 Usage

### Basic usage
```bash
./smm                    # Use default project
./smm -p PROJECT_ID      # Specify GCP project
```

## ⌨️ Keyboard Controls

### Navigation
| Key         | Action                                                     |
| ----------- | ---------------------------------------------------------- |
| `↑` `↓`     | Navigate list / Scroll in secret detail                   |
| `Tab`       | Switch focus between list and detail                      |
| `Shift + ←→`| Resize list view                                          |

### Search and Filtering
| Key         | Action                                                     |
| ----------- | ---------------------------------------------------------- |
| `/`         | Filter by secret name                                      |
| `Ctrl+F`    | Search in content of all secrets                          |

### Secret Management
| Key         | Action                                                     |
| ----------- | ---------------------------------------------------------- |
| `i`         | Show secret information (metadata, creation date, labels) |
| `c`         | Copy secret to clipboard                                   |
| `n`         | Create new version of secret                               |
| `v`         | Show/hide secret versions                                  |
| `r`         | Restore selected version                                   |

### System
| Key         | Action                                                     |
| ----------- | ---------------------------------------------------------- |
| `p`         | Switch GCP project                                         |
| `Esc`       | Refresh / Cancel operation                                 |
| `Ctrl+C`    | Exit program                                               |

### Command line options

| Option            | Description                                    |
| ----------------- | ---------------------------------------------- |
| `-p PROJECT_ID`   | Load secrets from specified project           |

## 🔐 Authentication

SMM uses existing `gcloud` authentication. Make sure you're authenticated before using the tool.

### Check authentication

```bash
gcloud config list
```

The output should contain your account and project:

```bash
[core]
account = your@email.com
project = your-gcp-project
```

### Check permissions

```bash
gcloud secrets list --project=your-gcp-project
```

If you have the necessary permissions, you'll see the project's secret list.

### Required permissions

Your account needs the following IAM roles:
- `roles/secretmanager.viewer` - To list and read secrets
- `roles/secretmanager.secretVersionManager` - To create new versions

### Authentication with gcloud
If you're not authenticated, you can do so with:
```bash
gcloud auth login
```

## 🎨 Syntax Highlighting

SMM automatically detects content format and applies syntax highlighting for:

- 🌱 **Bash/Env** - Environment variables
- 📄 **JSON** - Structured data  
- ⚙️ **INI** - Configuration files
- 🐘 **PHP** - PHP code

## 📁 Configuration

The application stores its configuration in `~/.config/smm/config.yaml`:

```yaml
projectIds: ["project-1", "project-2"]  # Available projects
selected: "project-1"                   # Selected project
logPath: "/path/to/log/file"            # Log file (optional)
```

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)  
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
