package client

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// LocalFileClient treats local files as secrets, scanning configured paths
// and managing them through the standard Client interface
type LocalFileClient struct {
	projectID string
	config    LocalClientConfig
}

// LocalClientConfig defines the configuration for the local file client
type LocalClientConfig struct {
	Paths []PathConfig `yaml:"paths"`
}

// PathConfig defines a path pattern with associated file types to scan
type PathConfig struct {
	Pattern   string   `yaml:"pattern"`    // Path pattern with wildcards, e.g., "/home/user/projects/*/config/"
	FileTypes []string `yaml:"file_types"` // File extensions or patterns, e.g., [".env", ".env.*", "*.key"]
	Name      string   `yaml:"name"`       // Optional name pattern for secrets, e.g., "[...] fileName" where ... represents directory levels
}

// NewLocalFileClient creates a new local file client with hardcoded configuration
// Following the plan, we start with hardcoded paths for initial implementation
func NewLocalFileClient(projectID string) (*LocalFileClient, error) {
	log.Info().Str("projectID", projectID).Msg("Creating LocalFileClient")

	// Hardcoded configuration as specified in the plan
	// Using paths within the SMM project for testing
	config := LocalClientConfig{
		Paths: []PathConfig{
			{
				Pattern:   "/home/felipe/Seafile/civitatis/tilt/tilts/*/",
				FileTypes: []string{".env", ".env.*"},
				Name:      "[.] fileName", // Example: [calendar-service] .env.tilt
			},
		},
	}

	client := &LocalFileClient{
		projectID: projectID,
		config:    config,
	}

	log.Info().Str("projectID", projectID).Msg("LocalFileClient initialized")
	return client, nil
}

// scanFiles discovers all files matching the configured patterns
func (l *LocalFileClient) scanFiles() ([]SecretInfo, error) {
	var allFiles []string

	for _, pathConfig := range l.config.Paths {
		files, err := l.expandPathPattern(pathConfig)
		if err != nil {
			log.Warn().Err(err).Str("pattern", pathConfig.Pattern).Msg("Failed to expand path pattern")
			continue
		}
		allFiles = append(allFiles, files...)
	}

	// Convert file paths to SecretInfo objects
	var secretInfos []SecretInfo
	for _, filePath := range allFiles {
		secretInfo, err := l.createSecretInfo(filePath)
		if err != nil {
			log.Warn().Err(err).Str("file", filePath).Msg("Failed to create SecretInfo for file")
			continue
		}
		secretInfos = append(secretInfos, secretInfo)
	}

	log.Debug().Int("count", len(secretInfos)).Msg("File scan completed")
	return secretInfos, nil
}

// fileMatchesPathConfig checks if a file path matches a specific PathConfig
func (l *LocalFileClient) fileMatchesPathConfig(filePath string, pathConfig PathConfig) bool {
	// Check if the file is in a directory that matches the pattern
	dirs := l.expandDirectoryPattern(pathConfig.Pattern)

	for _, dir := range dirs {
		if strings.HasPrefix(filePath, dir) {
			// Check if the file matches any of the file types
			fileName := filepath.Base(filePath)
			for _, fileType := range pathConfig.FileTypes {
				if l.fileMatchesType(fileName, fileType) {
					return true
				}
			}
		}
	}

	return false
}

// fileMatchesType checks if a filename matches a file type pattern
func (l *LocalFileClient) fileMatchesType(fileName, fileType string) bool {
	if strings.HasPrefix(fileType, ".") {
		// Extension patterns like ".env" or ".env.*"
		if strings.Contains(fileType, "*") {
			// For patterns like ".env.*", check if it starts with the base extension
			baseExt := strings.TrimSuffix(fileType, "*")
			return strings.Contains(fileName, baseExt)
		} else {
			// For exact extensions like ".env"
			return strings.HasSuffix(fileName, fileType)
		}
	} else {
		// Full filename patterns like "secrets.yaml" or "*.key"
		if strings.Contains(fileType, "*") {
			// Use filepath.Match for wildcard patterns
			matched, _ := filepath.Match(fileType, fileName)
			return matched
		} else {
			// Exact filename match
			return fileName == fileType
		}
	}
}

// generateSecretName generates a secret name based on the PathConfig.Name pattern
// If no pattern is specified, returns the original fileName
// Pattern format: "[...] fileName" where ... represents directory levels to traverse up
func (l *LocalFileClient) generateSecretName(filePath string, pathConfig PathConfig) string {
	fileName := filepath.Base(filePath)

	// If no name pattern is specified, return the original fileName for backward compatibility
	if pathConfig.Name == "" {
		return fileName
	}

	// Find the [...] pattern first
	startIdx := strings.Index(pathConfig.Name, "[")
	endIdx := strings.Index(pathConfig.Name, "]")

	// Check if pattern contains brackets
	if startIdx != -1 || endIdx != -1 {
		// If we have brackets, they must be well-formed
		if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
			// Malformed brackets - fallback to original filename
			return fileName
		}
	} else {
		// No brackets at all - treat as literal name pattern
		// Replace "fileName" placeholder with actual file name
		return strings.ReplaceAll(pathConfig.Name, "fileName", fileName)
	}

	// Extract the dots pattern
	dotsPattern := pathConfig.Name[startIdx+1 : endIdx]
	dotCount := strings.Count(dotsPattern, ".")

	if dotCount == 0 {
		return fileName // No dots means no directory traversal
	}

	// Navigate up the directory tree
	// Start from the directory containing the file
	currentPath := filepath.Dir(filePath)

	// For single dot [.], we want the directory name itself
	// For double dot [..], we want the parent directory name
	// For triple dot [...], we want the grandparent directory name
	targetPath := currentPath
	for i := 1; i < dotCount && targetPath != "/" && targetPath != "." && targetPath != ""; i++ {
		targetPath = filepath.Dir(targetPath)
	}

	// Get the directory name at the target level
	targetDirName := filepath.Base(targetPath)
	if targetDirName == "/" || targetDirName == "." || targetDirName == "" {
		targetDirName = "root" // Fallback name
	}

	// Replace the [...] pattern with [targetDirName]
	result := pathConfig.Name[:startIdx+1] + targetDirName + pathConfig.Name[endIdx:]

	// Replace fileName placeholder with actual file name
	result = strings.ReplaceAll(result, "fileName", fileName)

	return result
}

// expandPathPattern expands a PathConfig into actual file paths
func (l *LocalFileClient) expandPathPattern(pathConfig PathConfig) ([]string, error) {
	var allFiles []string

	// Handle patterns that may contain spaces by using a custom expansion
	dirs := l.expandDirectoryPattern(pathConfig.Pattern)

	// For each matching directory, find files matching the file type patterns
	for _, dir := range dirs {
		for _, fileType := range pathConfig.FileTypes {
			var pattern string

			// Handle different file type patterns
			if strings.HasPrefix(fileType, ".") {
				// Extension patterns like ".env" or ".env.*"
				if strings.Contains(fileType, "*") {
					pattern = filepath.Join(dir, "*"+fileType)
				} else {
					pattern = filepath.Join(dir, "*"+fileType)
				}
			} else {
				// Full filename patterns like "secrets.yaml" or "*.key"
				pattern = filepath.Join(dir, fileType)
			}

			files, err := filepath.Glob(pattern)
			if err != nil {
				log.Debug().Err(err).Str("pattern", pattern).Msg("Failed to expand file pattern")
				continue
			}

			// Filter out directories, we only want files
			for _, file := range files {
				if info, err := os.Stat(file); err == nil && !info.IsDir() {
					allFiles = append(allFiles, file)
				}
			}
		}
	}

	return l.deduplicateFiles(allFiles), nil
}

// expandDirectoryPattern handles directory patterns with wildcards, even with spaces in paths
func (l *LocalFileClient) expandDirectoryPattern(pattern string) []string {
	// First try the standard filepath.Glob
	if dirs, err := filepath.Glob(pattern); err == nil && len(dirs) > 0 {
		return dirs
	}

	// If that fails (usually due to spaces), try manual expansion
	// Split the pattern by the first wildcard
	wildcardPos := strings.Index(pattern, "*")
	if wildcardPos == -1 {
		// No wildcard, check if directory exists
		if info, err := os.Stat(pattern); err == nil && info.IsDir() {
			return []string{pattern}
		}
		return []string{}
	}

	// Find the directory part before the wildcard
	beforeWildcard := pattern[:wildcardPos]
	afterWildcard := pattern[wildcardPos+1:]

	// Find the last directory separator before the wildcard
	lastSep := strings.LastIndex(beforeWildcard, string(os.PathSeparator))
	if lastSep == -1 {
		return []string{}
	}

	baseDir := beforeWildcard[:lastSep]
	matchPrefix := beforeWildcard[lastSep+1:]

	// Read the base directory
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		log.Debug().Err(err).Str("baseDir", baseDir).Msg("Failed to read base directory")
		return []string{}
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), matchPrefix) {
			fullPath := filepath.Join(baseDir, entry.Name()) + afterWildcard
			// Verify the full path exists
			if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
				dirs = append(dirs, fullPath)
			}
		}
	}

	return dirs
}

// deduplicateFiles removes duplicate file paths
func (l *LocalFileClient) deduplicateFiles(files []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, file := range files {
		if !seen[file] {
			seen[file] = true
			result = append(result, file)
		}
	}

	return result
}

// createSecretInfo converts a file path into a SecretInfo object using the matching PathConfig
func (l *LocalFileClient) createSecretInfo(filePath string) (SecretInfo, error) {
	// Find the matching PathConfig for this file
	var matchingPathConfig *PathConfig
	for _, pathConfig := range l.config.Paths {
		if l.fileMatchesPathConfig(filePath, pathConfig) {
			matchingPathConfig = &pathConfig
			break
		}
	}

	return l.createSecretInfoWithConfig(filePath, matchingPathConfig)
}

// createSecretInfoWithConfig converts a file path into a SecretInfo object using a specific PathConfig
func (l *LocalFileClient) createSecretInfoWithConfig(filePath string, pathConfig *PathConfig) (SecretInfo, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return SecretInfo{}, fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	// Extract file name and directory information
	fileName := filepath.Base(filePath)
	dirPath := filepath.Dir(filePath)
	ext := filepath.Ext(filePath)

	// Create labels and annotations following the plan format
	labels := map[string]string{
		"type":      "local-file",
		"directory": dirPath,
		"extension": ext,
	}

	// Add specific labels based on file type
	if strings.Contains(fileName, ".env") {
		labels["file_type"] = "env-file"
	} else if strings.HasSuffix(fileName, ".json") {
		labels["file_type"] = "json-file"
	} else if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		labels["file_type"] = "yaml-file"
	} else if strings.HasSuffix(fileName, ".key") || strings.HasSuffix(fileName, ".pem") {
		labels["file_type"] = "key-file"
	}

	annotations := map[string]string{
		"size":          fmt.Sprintf("%d", fileInfo.Size()),
		"last_modified": fileInfo.ModTime().Format(time.RFC3339),
		"mode":          fileInfo.Mode().String(),
	}

	// Generate the secret name using the PathConfig if available
	var secretName string
	if pathConfig != nil {
		secretName = l.generateSecretName(filePath, *pathConfig)
	} else {
		secretName = fileName // Fallback to original behavior
	}

	return SecretInfo{
		Name:        secretName,
		FullPath:    filePath,
		CreateTime:  fileInfo.ModTime(), // Use modification time as creation time
		Labels:      labels,
		Annotations: annotations,
	}, nil
}

// Secrets returns all discovered files as SecretInfo objects
// Implements the Client interface
func (l *LocalFileClient) Secrets() ([]SecretInfo, error) {
	return l.scanFiles()
}

// GetSecret reads and returns the contents of the specified file
// As per the plan, this returns the current version of the file
func (l *LocalFileClient) GetSecret(secretName string) ([]byte, error) {
	filePath, err := l.findSecretPath(secretName)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	log.Debug().Str("secret", secretName).Str("path", filePath).
		Int("size", len(content)).Msg("Secret content retrieved")

	return content, nil
}

// GetSecretInfo returns detailed information about a specific secret/file
func (l *LocalFileClient) GetSecretInfo(fullPath string) (SecretInfo, error) {
	// If fullPath is actually a file path, return its info directly
	if _, err := os.Stat(fullPath); err == nil {
		return l.createSecretInfo(fullPath)
	}

	// Otherwise, search in discovered files
	secrets, err := l.scanFiles()
	if err != nil {
		return SecretInfo{}, fmt.Errorf("failed to scan files: %w", err)
	}

	for _, info := range secrets {
		if info.FullPath == fullPath || info.Name == fullPath {
			return info, nil
		}
	}

	return SecretInfo{}, fmt.Errorf("secret not found: %s", fullPath)
}

// GetSecretVersions returns version information for a secret
// As per the plan, local files only have one version (the current one)
func (l *LocalFileClient) GetSecretVersions(secretName string) ([]Version, error) {
	filePath, err := l.findSecretPath(secretName)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	// Local files only have one version - the current one
	version := Version{
		Name:      fmt.Sprintf("%s-current", secretName),
		State:     "enabled",
		Version:   1, // Always version 1 for local files
		FullPath:  filePath,
		CreatedAt: fileInfo.ModTime(),
	}

	return []Version{version}, nil
}

// GetSecretVersion returns the content of a specific version of a secret
// Since local files don't have historical versions, this always returns the current content
func (l *LocalFileClient) GetSecretVersion(secretName, version string) ([]byte, error) {
	// For local files, any version request returns the current content
	log.Debug().Str("secret", secretName).Str("version", version).
		Msg("Returning current version (local files don't support historical versions)")

	return l.GetSecret(secretName)
}

// CreateSecret creates a new empty file for the given secret name
func (l *LocalFileClient) CreateSecret(name string) error {
	if len(l.config.Paths) == 0 {
		return fmt.Errorf("no configured paths for LocalFileClient")
	}

	baseDir := strings.Replace(l.config.Paths[0].Pattern, "*", "", -1)
	baseDir = filepath.Clean(baseDir)
	fileExt := ".env"
	if len(l.config.Paths[0].FileTypes) > 0 {
		fileExt = strings.TrimPrefix(l.config.Paths[0].FileTypes[0], ".")
		fileExt = "." + strings.TrimSuffix(fileExt, ".*")
	}

	filePath := filepath.Join(baseDir, name+fileExt)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", baseDir, err)
	}

	emptyData := []byte{}
	if err := os.WriteFile(filePath, emptyData, 0644); err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	log.Info().Str("secret", name).Str("path", filePath).Msg("Created empty secret file")
	return nil
}

// AddSecretVersion creates a new version of a secret by overwriting the file
// Since local files don't support versioning, this overwrites the existing file
func (l *LocalFileClient) AddSecretVersion(secretName string, payload []byte) error {
	filePath, err := l.findSecretPath(secretName)
	if err != nil {
		return fmt.Errorf("secret not found for update: %w", err)
	}

	// Create backup directory if it doesn't exist
	backupDir := filepath.Dir(filePath) + "/.smm-backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		log.Warn().Err(err).Msg("Failed to create backup directory")
	} else {
		// Create a timestamped backup
		backupFile := filepath.Join(backupDir, fmt.Sprintf("%s.%d.backup",
			filepath.Base(filePath), time.Now().Unix()))
		if originalContent, err := os.ReadFile(filePath); err == nil {
			if err := os.WriteFile(backupFile, originalContent, 0644); err != nil {
				log.Warn().Err(err).Str("backup", backupFile).Msg("Failed to create backup")
			} else {
				log.Info().Str("backup", backupFile).Msg("Created backup before overwrite")
			}
		}
	}

	// Write the new content
	if err := os.WriteFile(filePath, payload, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	log.Info().Str("secret", secretName).Str("path", filePath).
		Int("size", len(payload)).Msg("Secret version added (file overwritten)")

	return nil
}

// SearchInSecrets searches for the given query within the content of all files
func (l *LocalFileClient) SearchInSecrets(query string) ([]SecretInfo, error) {
	secrets, err := l.scanFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to scan files: %w", err)
	}

	var matchingSecrets []SecretInfo
	query = strings.ToLower(query)

	for _, secretInfo := range secrets {
		// Search in filename first
		if strings.Contains(strings.ToLower(secretInfo.Name), query) {
			matchingSecrets = append(matchingSecrets, secretInfo)
			continue
		}

		// Search in file content
		content, err := os.ReadFile(secretInfo.FullPath)
		if err != nil {
			log.Debug().Err(err).Str("file", secretInfo.FullPath).Msg("Failed to read file for search")
			continue
		}

		if strings.Contains(strings.ToLower(string(content)), query) {
			matchingSecrets = append(matchingSecrets, secretInfo)
		}
	}

	log.Debug().Str("query", query).Int("matches", len(matchingSecrets)).Msg("Search completed")
	return matchingSecrets, nil
}

// findSecretPath finds the full file path for a given secret name
func (l *LocalFileClient) DeleteSecret(name string) error {
	filePath, err := l.findSecretPath(name)
	if err != nil {
		return fmt.Errorf("secret not found for deletion: %w", err)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", filePath, err)
	}

	log.Info().Str("secret", name).Str("path", filePath).Msg("Secret file deleted")
	return nil
}

func (l *LocalFileClient) DestroySecretVersion(name string, version int) error {
	return fmt.Errorf("version deletion not supported for local file client")
}

func (l *LocalFileClient) findSecretPath(secretName string) (string, error) {
	secrets, err := l.scanFiles()
	if err != nil {
		return "", fmt.Errorf("failed to scan files: %w", err)
	}

	// Look for exact name match first
	for _, info := range secrets {
		if info.FullPath == secretName {
			return info.FullPath, nil
		}
	}

	// Look for partial name match (useful for files with extensions)
	for _, info := range secrets {
		if strings.Contains(info.Name, secretName) {
			return info.FullPath, nil
		}
	}

	return "", fmt.Errorf("secret not found: %s", secretName)
}
