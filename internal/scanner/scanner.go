package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// AudioExtensions are common audio file extensions
var AudioExtensions = []string{
	".wav", ".mp3", ".flac", ".aif", ".aiff", ".ogg", ".m4a", ".wma", ".aac",
}

// SampleFile represents a discovered audio sample file
type SampleFile struct {
	OriginalPath string
	FileName     string
	Extension    string
}

// ScanDirectory recursively scans a directory for audio files
func ScanDirectory(dir string) ([]SampleFile, error) {
	var samples []SampleFile

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		for _, audioExt := range AudioExtensions {
			if ext == audioExt {
				samples = append(samples, SampleFile{
					OriginalPath: path,
					FileName:     filepath.Base(path),
					Extension:    ext,
				})
				break
			}
		}

		return nil
	})

	return samples, err
}
