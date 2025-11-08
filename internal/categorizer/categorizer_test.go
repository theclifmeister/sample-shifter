package categorizer

import (
	"testing"

	"github.com/theclifmeister/sample-shifter/internal/scanner"
)

func TestCategorize(t *testing.T) {
	tests := []struct {
		fileName         string
		extension        string
		expectedCategory Category
	}{
		{"kick_01.wav", ".wav", CategoryDrum},
		{"snare_heavy.wav", ".wav", CategoryDrum},
		{"hihat_closed.wav", ".wav", CategoryDrum},
		{"clap_01.wav", ".wav", CategoryDrum},
		{"bass_sub.wav", ".wav", CategoryBass},
		{"808_bass.wav", ".wav", CategoryBass},
		{"synth_lead.wav", ".wav", CategorySynth},
		{"synth_pad.mp3", ".mp3", CategorySynth},
		{"vocal_shot.wav", ".wav", CategoryOneShot}, // "shot" keyword has higher priority (OneShot category)
		{"vocal_sample.wav", ".wav", CategoryVocal},
		{"fx_riser.wav", ".wav", CategoryFX},
		{"shaker_loop.wav", ".wav", CategoryPercussion},
		{"piano_chord.wav", ".wav", CategoryMelodic},
		{"guitar_strum.wav", ".wav", CategoryMelodic},
		{"bell_ring.flac", ".flac", CategoryMelodic},
		{"strings_melody.wav", ".wav", CategoryMelodic},
		{"beat_loop.wav", ".wav", CategoryLoop},
		{"loop_128bpm.wav", ".wav", CategoryLoop},
		{"oneshot_stab.wav", ".wav", CategoryOneShot},
		{"random_sound.wav", ".wav", CategoryUncategorized},
	}

	targetDir := "/tmp/test-target"

	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			sample := scanner.SampleFile{
				OriginalPath: "/test/path/" + tt.fileName,
				FileName:     tt.fileName,
				Extension:    tt.extension,
			}

			result := Categorize(sample, targetDir, false)

			if result.Category != tt.expectedCategory {
				t.Errorf("Expected category %s for %s, got %s",
					tt.expectedCategory, tt.fileName, result.Category)
			}

			if result.Sample.FileName != tt.fileName {
				t.Errorf("Expected filename %s, got %s", tt.fileName, result.Sample.FileName)
			}

			if result.TargetPath == "" {
				t.Error("TargetPath should not be empty")
			}
		})
	}
}

func TestCategorizeBatch(t *testing.T) {
	samples := []scanner.SampleFile{
		{OriginalPath: "/test/kick.wav", FileName: "kick.wav", Extension: ".wav"},
		{OriginalPath: "/test/bass.wav", FileName: "bass.wav", Extension: ".wav"},
		{OriginalPath: "/test/synth.wav", FileName: "synth.wav", Extension: ".wav"},
	}

	targetDir := "/tmp/test-target"
	categorized := CategorizeBatch(samples, targetDir, false)

	if len(categorized) != len(samples) {
		t.Errorf("Expected %d categorized files, got %d", len(samples), len(categorized))
	}

	// Check that each file was categorized
	for i, cat := range categorized {
		if cat.Sample.FileName != samples[i].FileName {
			t.Errorf("Sample order mismatch at index %d", i)
		}
		if cat.Category == "" {
			t.Errorf("Category should not be empty for %s", cat.Sample.FileName)
		}
	}
}

func TestCategorizeCaseInsensitive(t *testing.T) {
	tests := []struct {
		fileName string
		expected Category
	}{
		{"KICK_01.wav", CategoryDrum},
		{"Snare_Heavy.wav", CategoryDrum},
		{"BASS_SUB.wav", CategoryBass},
		{"Synth_Lead.wav", CategorySynth},
	}

	targetDir := "/tmp/test-target"

	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			sample := scanner.SampleFile{
				OriginalPath: "/test/" + tt.fileName,
				FileName:     tt.fileName,
				Extension:    ".wav",
			}

			result := Categorize(sample, targetDir, false)

			if result.Category != tt.expected {
				t.Errorf("Case-insensitive matching failed: expected %s for %s, got %s",
					tt.expected, tt.fileName, result.Category)
			}
		})
	}
}

func TestCategoryKeywords(t *testing.T) {
	// Verify that each category has keywords defined
	if len(keywords) == 0 {
		t.Error("Keywords map should not be empty")
	}

	for category, words := range keywords {
		if len(words) == 0 {
			t.Errorf("Category %s has no keywords", category)
		}
	}
}
