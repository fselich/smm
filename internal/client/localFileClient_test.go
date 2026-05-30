package client

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLocalFileClient(t *testing.T) {
	client, err := NewLocalFileClient("test-local-project")
	if err != nil {
		t.Fatalf("Failed to create LocalFileClient: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.projectID != "test-local-project" {
		t.Errorf("Expected projectID to be 'test-local-project', got '%s'", client.projectID)
	}

	if len(client.config.Paths) == 0 {
		t.Error("Expected config paths to be populated")
	}

	// fileCache no longer exists in simplified implementation
}

func TestLocalFileClient_Interface(t *testing.T) {
	// Verify that LocalFileClient implements the Client interface
	var _ Client = &LocalFileClient{}
}

func TestLocalFileClient_Secrets(t *testing.T) {
	client, err := NewLocalFileClient("test-local-project")
	if err != nil {
		t.Fatalf("Failed to create LocalFileClient: %v", err)
	}

	secrets, err := client.Secrets()
	if err != nil {
		t.Errorf("Secrets() returned an error: %v", err)
	}

	// Secrets should be a slice (could be empty if no matching files exist)
	// Note: secrets is never nil, it's initialized as empty slice

	t.Logf("Found %d secrets", len(secrets))
	for i, secret := range secrets {
		if i >= 5 { // Only log first 5 to avoid spam
			t.Logf("... and %d more", len(secrets)-i)
			break
		}
		t.Logf("Secret %d: %s (%s)", i+1, secret.Name, secret.FullPath)
	}
}

func TestLocalFileClient_CreateSecretInfo(t *testing.T) {
	client, err := NewLocalFileClient("test-project")
	if err != nil {
		t.Fatalf("Failed to create LocalFileClient: %v", err)
	}

	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.env")
	testContent := "DATABASE_URL=localhost:5432\nAPI_KEY=secret123"

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	secretInfo, err := client.createSecretInfo(testFile)
	if err != nil {
		t.Fatalf("createSecretInfo failed: %v", err)
	}

	if secretInfo.Name != "test.env" {
		t.Errorf("Expected name 'test.env', got '%s'", secretInfo.Name)
	}

	if secretInfo.FullPath != testFile {
		t.Errorf("Expected FullPath '%s', got '%s'", testFile, secretInfo.FullPath)
	}

	if secretInfo.Labels["type"] != "local-file" {
		t.Errorf("Expected type label 'local-file', got '%s'", secretInfo.Labels["type"])
	}

	if secretInfo.Labels["extension"] != ".env" {
		t.Errorf("Expected extension label '.env', got '%s'", secretInfo.Labels["extension"])
	}

	if secretInfo.Labels["file_type"] != "env-file" {
		t.Errorf("Expected file_type label 'env-file', got '%s'", secretInfo.Labels["file_type"])
	}
}

func TestLocalFileClient_GetSecretVersions(t *testing.T) {
	client, err := NewLocalFileClient("test-project")
	if err != nil {
		t.Fatalf("Failed to create LocalFileClient: %v", err)
	}

	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.env")
	testContent := "DATABASE_URL=localhost:5432"

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// In the simplified implementation, we need to modify the client config to scan the temp directory
	client.config = LocalClientConfig{
		Paths: []PathConfig{
			{
				Pattern:   tmpDir + "/",
				FileTypes: []string{".env"},
			},
		},
	}

	versions, err := client.GetSecretVersions("test.env")
	if err != nil {
		t.Fatalf("GetSecretVersions failed: %v", err)
	}

	if len(versions) != 1 {
		t.Errorf("Expected 1 version, got %d", len(versions))
	}

	version := versions[0]
	if version.Version != 1 {
		t.Errorf("Expected version 1, got %d", version.Version)
	}

	if version.State != "enabled" {
		t.Errorf("Expected state 'enabled', got '%s'", version.State)
	}

	if version.FullPath != testFile {
		t.Errorf("Expected FullPath '%s', got '%s'", testFile, version.FullPath)
	}
}

func TestLocalFileClient_GetSecret(t *testing.T) {
	client, err := NewLocalFileClient("test-project")
	if err != nil {
		t.Fatalf("Failed to create LocalFileClient: %v", err)
	}

	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.env")
	testContent := "DATABASE_URL=localhost:5432\nAPI_KEY=secret123"

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// In the simplified implementation, we need to modify the client config to scan the temp directory
	client.config = LocalClientConfig{
		Paths: []PathConfig{
			{
				Pattern:   tmpDir + "/",
				FileTypes: []string{".env"},
			},
		},
	}

	content, err := client.GetSecret("test.env")
	if err != nil {
		t.Fatalf("GetSecret failed: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content '%s', got '%s'", testContent, string(content))
	}
}

func TestLocalFileClient_SearchInSecrets(t *testing.T) {
	client, err := NewLocalFileClient("test-project")
	if err != nil {
		t.Fatalf("Failed to create LocalFileClient: %v", err)
	}

	// Create temporary files for testing
	tmpDir := t.TempDir()

	// File 1: Contains "database"
	testFile1 := filepath.Join(tmpDir, "db.env")
	testContent1 := "DATABASE_URL=localhost:5432"
	err = os.WriteFile(testFile1, []byte(testContent1), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}

	// File 2: Contains "api"
	testFile2 := filepath.Join(tmpDir, "api.env")
	testContent2 := "API_KEY=secret123"
	err = os.WriteFile(testFile2, []byte(testContent2), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}

	// In the simplified implementation, we need to modify the client config to scan the temp directory
	client.config = LocalClientConfig{
		Paths: []PathConfig{
			{
				Pattern:   tmpDir + "/",
				FileTypes: []string{".env"},
			},
		},
	}

	// Search for "database" - should find file 1
	results, err := client.SearchInSecrets("database")
	if err != nil {
		t.Fatalf("SearchInSecrets failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'database', got %d", len(results))
	} else if results[0].Name != "db.env" {
		t.Errorf("Expected result 'db.env', got '%s'", results[0].Name)
	}

	// Search for "api" - should find file 2 by both filename and content
	results, err = client.SearchInSecrets("api")
	if err != nil {
		t.Fatalf("SearchInSecrets failed: %v", err)
	}

	if len(results) != 1 { // Should find api.env (matches filename, so doesn't check content)
		t.Errorf("Expected 1 result for 'api', got %d", len(results))
	}
}

// TestLocalFileClient_GenerateSecretName tests the new secret naming functionality
func TestLocalFileClient_GenerateSecretName(t *testing.T) {
	client := &LocalFileClient{}

	tests := []struct {
		name         string
		filePath     string
		pathConfig   PathConfig
		expectedName string
		description  string
	}{
		{
			name:         "No name pattern - backward compatibility",
			filePath:     "/home/user/projects/calendar-service/config/.env.dev",
			pathConfig:   PathConfig{Name: ""},
			expectedName: ".env.dev",
			description:  "Should return original filename when no Name pattern is specified",
		},
		{
			name:         "Simple literal pattern",
			filePath:     "/home/user/projects/calendar-service/config/.env.dev",
			pathConfig:   PathConfig{Name: "service-fileName"},
			expectedName: "service-.env.dev",
			description:  "Should replace fileName placeholder with actual filename",
		},
		{
			name:         "Single dot pattern",
			filePath:     "/home/user/projects/calendar-service/config/.env.dev",
			pathConfig:   PathConfig{Name: "[.] fileName"},
			expectedName: "[config] .env.dev",
			description:  "Single dot should go up one directory level",
		},
		{
			name:         "Double dot pattern",
			filePath:     "/home/user/projects/calendar-service/config/.env.dev",
			pathConfig:   PathConfig{Name: "[..] fileName"},
			expectedName: "[calendar-service] .env.dev",
			description:  "Double dot should go up two directory levels",
		},
		{
			name:         "Triple dot pattern",
			filePath:     "/home/user/projects/calendar-service/config/.env.dev",
			pathConfig:   PathConfig{Name: "[...] fileName"},
			expectedName: "[projects] .env.dev",
			description:  "Triple dot should go up three directory levels",
		},
		{
			name:         "Root directory edge case",
			filePath:     "/config/.env.dev",
			pathConfig:   PathConfig{Name: "[...] fileName"},
			expectedName: "[root] .env.dev",
			description:  "Should fallback to 'root' when reaching filesystem root",
		},
		{
			name:         "Complex pattern with prefix and suffix",
			filePath:     "/home/user/projects/calendar-service/config/.env.dev",
			pathConfig:   PathConfig{Name: "prod-[..]-fileName"},
			expectedName: "prod-[calendar-service]-.env.dev",
			description:  "Should handle complex patterns with prefix and suffix",
		},
		{
			name:         "Malformed pattern - no closing bracket",
			filePath:     "/home/user/projects/calendar-service/config/.env.dev",
			pathConfig:   PathConfig{Name: "[.. fileName"},
			expectedName: ".env.dev",
			description:  "Should fallback to filename for malformed patterns",
		},
		{
			name:         "Empty dots pattern",
			filePath:     "/home/user/projects/calendar-service/config/.env.dev",
			pathConfig:   PathConfig{Name: "[] fileName"},
			expectedName: ".env.dev",
			description:  "Should fallback to filename when no dots in pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.generateSecretName(tt.filePath, tt.pathConfig)
			if result != tt.expectedName {
				t.Errorf("generateSecretName() = %q, want %q\nDescription: %s", result, tt.expectedName, tt.description)
			}
		})
	}
}

// TestLocalFileClient_FileMatchesType tests the file type matching functionality
func TestLocalFileClient_FileMatchesType(t *testing.T) {
	client := &LocalFileClient{}

	tests := []struct {
		name     string
		fileName string
		fileType string
		expected bool
	}{
		{"Exact extension match", ".env.dev", ".env", false}, // .env should match files ending with .env
		{"Extension with wildcard", ".env.dev", ".env.*", true},
		{"Exact filename match", "secrets.yaml", "secrets.yaml", true},
		{"Wildcard filename match", "secrets.yaml", "*.yaml", true},
		{"No match", "config.json", ".env", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.fileMatchesType(tt.fileName, tt.fileType)
			if result != tt.expected {
				t.Errorf("fileMatchesType(%q, %q) = %t, want %t", tt.fileName, tt.fileType, result, tt.expected)
			}
		})
	}
}

// TestLocalFileClient_CreateSecretInfoWithNewName tests the complete integration
func TestLocalFileClient_CreateSecretInfoWithNewName(t *testing.T) {
	client, err := NewLocalFileClient("test-project")
	if err != nil {
		t.Fatalf("Failed to create LocalFileClient: %v", err)
	}

	// Create a temporary directory structure
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "projects", "calendar-service", "config")
	err = os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	testFile := filepath.Join(projectDir, ".env.dev")
	testContent := "DATABASE_URL=localhost:5432\nAPI_KEY=secret123"

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with name pattern
	pathConfigWithName := PathConfig{
		Pattern:   filepath.Join(tmpDir, "projects/*/config/"),
		FileTypes: []string{".env.*"},
		Name:      "[..] fileName",
	}

	secretInfo, err := client.createSecretInfoWithConfig(testFile, &pathConfigWithName)
	if err != nil {
		t.Fatalf("createSecretInfoWithConfig failed: %v", err)
	}

	expectedName := "[calendar-service] .env.dev"
	if secretInfo.Name != expectedName {
		t.Errorf("Expected name %q, got %q", expectedName, secretInfo.Name)
	}

	// Test backward compatibility (no name pattern)
	pathConfigNoName := PathConfig{
		Pattern:   filepath.Join(tmpDir, "projects/*/config/"),
		FileTypes: []string{".env.*"},
		// Name field intentionally omitted
	}

	secretInfo2, err := client.createSecretInfoWithConfig(testFile, &pathConfigNoName)
	if err != nil {
		t.Fatalf("createSecretInfoWithConfig failed: %v", err)
	}

	if secretInfo2.Name != ".env.dev" {
		t.Errorf("Expected backward compatibility name '.env.dev', got %q", secretInfo2.Name)
	}
}

// TestLocalFileClient_shouldRefreshCache is no longer relevant as caching was removed
// The simplified implementation always scans fresh
