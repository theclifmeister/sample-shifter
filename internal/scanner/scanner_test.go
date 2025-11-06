package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDirectory(t *testing.T) {
	// Create a temporary test directory
	tmpDir := t.TempDir()

	// Create test files
	testFiles := []struct {
		name      string
		isAudio   bool
	}{
		{"kick.wav", true},
		{"snare.mp3", true},
		{"bass.flac", true},
		{"readme.txt", false},
		{"image.jpg", false},
		{"synth.aif", true},
	}

	for _, tf := range testFiles {
		path := filepath.Join(tmpDir, tf.name)
		if err := os.WriteFile(path, []byte{}, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Create a subdirectory with more files
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	
	subFile := filepath.Join(subDir, "vocal.ogg")
	if err := os.WriteFile(subFile, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create subdirectory file: %v", err)
	}

	// Scan the directory
	samples, err := ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	// We should find 5 audio files (4 in root + 1 in subdir)
	expectedCount := 5
	if len(samples) != expectedCount {
		t.Errorf("Expected %d audio files, got %d", expectedCount, len(samples))
	}

	// Verify each sample has required fields
	for _, sample := range samples {
		if sample.OriginalPath == "" {
			t.Error("Sample has empty OriginalPath")
		}
		if sample.FileName == "" {
			t.Error("Sample has empty FileName")
		}
		if sample.Extension == "" {
			t.Error("Sample has empty Extension")
		}
	}
}

func TestScanEmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	samples, err := ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if len(samples) != 0 {
		t.Errorf("Expected 0 samples in empty directory, got %d", len(samples))
	}
}

func TestScanNonexistentDirectory(t *testing.T) {
	_, err := ScanDirectory("/path/that/does/not/exist")
	if err == nil {
		t.Error("Expected error for nonexistent directory, got nil")
	}
}

func TestAudioExtensions(t *testing.T) {
	// Verify that common audio extensions are included
	expectedExts := []string{".wav", ".mp3", ".flac", ".aif", ".aiff"}
	
	for _, ext := range expectedExts {
		found := false
		for _, audioExt := range AudioExtensions {
			if ext == audioExt {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected extension %s not found in AudioExtensions", ext)
		}
	}
}
