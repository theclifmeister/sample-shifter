package categorizer

import (
	"path/filepath"
	"strings"

	"github.com/theclifmeister/sample-shifter/internal/scanner"
)

// Category represents a sample category
type Category string

const (
	CategoryDrum          Category = "drums"
	CategoryBass          Category = "bass"
	CategorySynth         Category = "synth"
	CategoryVocal         Category = "vocals"
	CategoryFX            Category = "fx"
	CategoryPercussion    Category = "percussion"
	CategoryMelodic       Category = "melodic"
	CategoryLoop          Category = "loops"
	CategoryOneShot       Category = "oneshots"
	CategoryUncategorized Category = "uncategorized"
)

// CategorizedFile represents a file with its determined category
type CategorizedFile struct {
	Sample     scanner.SampleFile
	Category   Category
	TargetPath string
}

// keywords maps category keywords to categories
var keywords = map[Category][]string{
	CategoryDrum: {
		"kick", "snare", "hihat", "hi-hat", "clap", "tom", "cymbal", "crash", "ride",
		"drum", "bd", "sd", "hh", "ch", "oh",
	},
	CategoryBass: {
		"bass", "sub", "808", "909",
	},
	CategorySynth: {
		"synth", "lead", "pad", "pluck", "saw", "square", "sine",
	},
	CategoryVocal: {
		"vocal", "vox", "voice", "acapella", "choir", "shout", "chant",
	},
	CategoryFX: {
		"fx", "sfx", "riser", "downsweep", "whoosh", "impact", "sweep", "noise",
		"white", "reverse", "rev",
	},
	CategoryPercussion: {
		"perc", "percussion", "shaker", "conga", "bongo", "tambourine", "cowbell",
	},
	CategoryLoop: {
		"loop", "phrase", "bar", "beat",
	},
	CategoryOneShot: {
		"oneshot", "one-shot", "hit", "stab",
	},
}

// Categorize determines the category of a sample file based on its name
func Categorize(sample scanner.SampleFile, targetDir string) CategorizedFile {
	fileName := strings.ToLower(sample.FileName)
	nameWithoutExt := strings.ToLower(strings.TrimSuffix(sample.FileName, sample.Extension))

	category := CategoryUncategorized

	// Check if filename contains any keywords
	for cat, words := range keywords {
		for _, word := range words {
			if strings.Contains(fileName, word) || strings.Contains(nameWithoutExt, word) {
				category = cat
				goto categorized
			}
		}
	}

categorized:
	targetPath := filepath.Join(targetDir, string(category), sample.FileName)

	return CategorizedFile{
		Sample:     sample,
		Category:   category,
		TargetPath: targetPath,
	}
}

// CategorizeBatch categorizes multiple sample files
func CategorizeBatch(samples []scanner.SampleFile, targetDir string) []CategorizedFile {
	categorized := make([]CategorizedFile, 0, len(samples))

	for _, sample := range samples {
		categorized = append(categorized, Categorize(sample, targetDir))
	}

	return categorized
}
