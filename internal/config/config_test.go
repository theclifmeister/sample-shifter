package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigDefault(t *testing.T) {
	// Test that default config is returned when no path is provided
	config, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig with empty path should not error: %v", err)
	}

	if config == nil {
		t.Fatal("Config should not be nil")
	}

	if len(config.Categories) == 0 {
		t.Fatal("Default config should have categories")
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	testConfig := &CategoryConfig{
		Categories: []CategoryDefinition{
			{
				Name:     "test_category",
				Priority: 1,
				Keywords: []string{"test", "sample"},
				Subcategories: map[string][]string{
					"sub1": {"keyword1", "keyword2"},
				},
			},
		},
	}

	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load the config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig should not error: %v", err)
	}

	if len(config.Categories) != 1 {
		t.Errorf("Expected 1 category, got %d", len(config.Categories))
	}

	if config.Categories[0].Name != "test_category" {
		t.Errorf("Expected category name 'test_category', got '%s'", config.Categories[0].Name)
	}

	if len(config.Categories[0].Keywords) != 2 {
		t.Errorf("Expected 2 keywords, got %d", len(config.Categories[0].Keywords))
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	// Create a temporary invalid JSON file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(configPath, []byte("invalid json{"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("LoadConfig should error on invalid JSON")
	}
}

func TestLoadConfigNonexistent(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("LoadConfig should error on nonexistent file")
	}
}

func TestValidateConfigEmpty(t *testing.T) {
	config := &CategoryConfig{
		Categories: []CategoryDefinition{},
	}

	err := validateConfig(config)
	if err == nil {
		t.Error("Validation should fail for empty categories")
	}
}

func TestValidateConfigNoName(t *testing.T) {
	config := &CategoryConfig{
		Categories: []CategoryDefinition{
			{
				Name:     "",
				Priority: 1,
				Keywords: []string{"test"},
			},
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Error("Validation should fail for empty category name")
	}
}

func TestValidateConfigDuplicateName(t *testing.T) {
	config := &CategoryConfig{
		Categories: []CategoryDefinition{
			{
				Name:     "test",
				Priority: 1,
				Keywords: []string{"test1"},
			},
			{
				Name:     "test",
				Priority: 2,
				Keywords: []string{"test2"},
			},
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Error("Validation should fail for duplicate category names")
	}
}

func TestValidateConfigNoKeywords(t *testing.T) {
	config := &CategoryConfig{
		Categories: []CategoryDefinition{
			{
				Name:     "test",
				Priority: 1,
				Keywords: []string{},
			},
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Error("Validation should fail for category with no keywords")
	}
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()

	if config == nil {
		t.Fatal("Default config should not be nil")
	}

	if len(config.Categories) == 0 {
		t.Fatal("Default config should have categories")
	}

	// Verify some expected categories exist
	expectedCategories := []string{"drums", "bass", "synth", "vocals", "fx", "percussion", "melodic", "loops", "oneshots"}
	foundCategories := make(map[string]bool)

	for _, cat := range config.Categories {
		foundCategories[cat.Name] = true
	}

	for _, expected := range expectedCategories {
		if !foundCategories[expected] {
			t.Errorf("Expected category '%s' not found in default config", expected)
		}
	}
}
