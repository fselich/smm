# Changelog

## [0.2.5] - 2026-05-30

### Added
- Interactive secret creation form (`ctrl+n`)
- Brace expansion for creating multiple secrets at once (e.g. `secret-{dev,staging,prod}`)
- Migration from bubbletea v1 to bubbletea v2

### Changed
- Refactored model/view separation for secrets page
- Updated list delegate for bubbletea v2 compatibility
- Simplified overlay rendering

### Fixed
- Tests updated for bubbletea v2 compatibility

### Security
- Updated golang.org/x/crypto to v0.45.0
- Updated google.golang.org/grpc to v1.79.3
- Updated go.opentelemetry.io/otel to v1.41.0

## [0.1.14] - 2025-09-12

### Fixed
- Remove trailing newlines when editing secrets to prevent breaking microservices (resolves #12)

## [0.1.12] - 2025-08-16

### Changed
- Use context with cancel in GCP instance for better resource management

### Fixed  
- Improved security in written files and directories

## [0.1.11] - 2025-08-12

### Added
- Version flag (-v) to display version

## [0.1.7] - 2025-08-12

### Changed
- Replaced clipboard dependency for better cross-platform compatibility

## [0.1.6] - 2025-08-11

### Changed
- Replaced viper dependency with internal config package
- Refactored configuration management for better maintainability

## [0.1.5] - 2025-08-11

### Changed
- Build target limited to Linux only

## [0.1.4] - 2025-08-10

### Improved
- Enhanced secret info modal with better layout and information display
- Updated documentation (README files)

## [0.1.3] - 2025-08-09

### Added
- Secret info modal accessible with 'i' key showing metadata (name, path, creation date, labels, annotations)
- Rich SecretInfo struct replacing simple string lists for better metadata handling
- Real creation dates and metadata display instead of hardcoded fake data

### Improved
- Enhanced client interface with proper error handling for all methods
- Better API consistency between Version and SecretInfo structs
- More efficient caching of secret metadata from GCP API

### Fixed
- Updated dependencies: golang.org/x/oauth2 to v0.27.0, golang.org/x/crypto to v0.36.0, golang.org/x/net to v0.38.0

## [0.1.2] - 2025-08-08

### Fixed
- Fixed infinite recursion bug in Model.View() causing crashes
- Optimized string concatenation in overlay rendering for better performance

## [0.1.1] - 2025-08-07

### Added
- Test coverage for editor
- First public version
