**English** | **[Español](README.es.md)**

# SMM - Secret Manager Manager

**SMM** is a Terminal User Interface (TUI) tool that allows you to efficiently and securely view, edit, and manage Google Cloud Platform secrets.

## Features

- **Intuitive navigation** with modern terminal interface
- **Secret editing** with your favorite editor
- **Version management** - view, restore, and create new versions
- **Advanced search** - search by name or content
- **Copy to clipboard** with a single command
- **Syntax highlighting** for multiple formats (JSON, Bash, INI, PHP)
- **Multi-project** - easily switch between GCP projects

## Installation

### Homebrew (Recommended)

```bash
brew tap fselich/tap
brew install smm
```

### Download Binary

Download the latest release for your platform from [GitHub Releases](https://github.com/fselich/smm/releases):

#### Linux x64
```bash
wget https://github.com/fselich/smm/releases/latest/download/smm-linux-amd64.tar.gz
tar -xzf smm-linux-amd64.tar.gz
chmod +x smm
sudo mv smm /usr/local/bin/  # Optional: install system-wide
```

#### Linux ARM64
```bash
wget https://github.com/fselich/smm/releases/latest/download/smm-linux-arm64.tar.gz
tar -xzf smm-linux-arm64.tar.gz
chmod +x smm
sudo mv smm /usr/local/bin/  # Optional: install system-wide
```

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

## Usage

### Basic usage
```bash
./smm                    # Use default project
./smm -p PROJECT_ID      # Specify GCP project
./smm -v                 # Show version information
```

## Keyboard Controls

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
| `-v`              | Show version information                       |

## Authentication

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

## Syntax Highlighting

SMM automatically detects content format and applies syntax highlighting for:

- **Bash/Env** - Environment variables
- **JSON** - Structured data  
- **INI** - Configuration files
- **PHP** - PHP code

## Configuration

The application stores its configuration in `~/.config/smm/config.yaml`. **The configuration file is created automatically if it doesn't exist** when you first run SMM.

Configuration structure:

```yaml
projects:                               # List of GCP projects
  - id: "my-gcp-project-1"
    type: "gcp"
  - id: "my-gcp-project-2"  
    type: "gcp"
selected: "my-gcp-project-1"            # Currently selected project
logPath: "/path/to/log/file"            # Log file path (optional)
```

**Notes:**
- Projects are automatically added when you switch to them using the `p` key
- The `selected` field remembers your last used project  
- `logPath` is optional - leave empty to disable logging

## Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)  
5. Open a Pull Request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
