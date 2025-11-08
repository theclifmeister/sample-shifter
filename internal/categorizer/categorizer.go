package categorizer

import (
	"path/filepath"
	"regexp"
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
	CategoryAmbiance      Category = "ambiance"
	CategoryTransition    Category = "transition"
	CategoryFoley         Category = "foley"
	CategoryUncategorized Category = "uncategorized"
)

// CategorizedFile represents a file with its determined category
type CategorizedFile struct {
	Sample      scanner.SampleFile
	Category    Category
	Subcategory string
	TargetPath  string
}

// categoryPriority defines the order in which categories are checked
// Categories earlier in the list have higher priority
var categoryPriority = []Category{
	CategoryOneShot,     // Check one-shots first (shots/hits/stabs)
	CategoryDrum,
	CategoryBass,
	CategoryPercussion,
	CategoryVocal,
	CategorySynth,
	CategoryMelodic,
	CategoryFX,
	CategoryTransition,
	CategoryAmbiance,
	CategoryFoley,
	CategoryLoop,
}

// keywords maps category keywords to categories
var keywords = map[Category][]string{
	CategoryDrum: {
		"kick", "snare", "hihat", "hi-hat", "hi_hat", "hi hat", "hats", "clap", "tom", "cymbal", "crash", "ride",
		"drum", "bd", "sd", "hh", "closed hat", "open hat", "hat closed", "hat open",
		"sidestick", "side stick", "rimshot", "rim shot", "cup", "rim", "cym", "china",
		"crossstick", "cross stick",
	},
	CategoryBass: {
		"bass", "sub", "808", "909",
	},
	CategorySynth: {
		"synth", "lead", "pad", "pluck", "saw", "square", "sine",
	},
	CategoryVocal: {
		"vocal", "vox", "voice", "acapella", "choir", "shout", "chant", "adlib",
	},
	CategoryFX: {
		"fx", "sfx", "riser", "downsweep", "whoosh", "impact", "sweep", "noise",
		"white", "reverse", "rev", "glitch", "tone", "envelope", "pulse",
		"ufo", "bleeps", "sync", "click",
	},
	CategoryPercussion: {
		"perc", "percussion", "shaker", "conga", "bongo", "tambourine", "tamb", "cowbell",
		"cabasa", "clave", "claves", "agogo", "timbale", "timpani", "maracas", "maraca",
		"woodblock", "wood block", "triangle", "guiro", "djembe", "udu",
		"brush", "chk", "cowb",
	},
	CategoryMelodic: {
		"piano", "guitar", "bell", "marimba", "xylophone", "harp", "strings",
		"violin", "cello", "flute", "horn", "trumpet", "sax", "saxophone",
		"organ", "keys", "brass", "woodwind", "arpeggio", "arpeggiated", "melody",
		"oud", "bouzouki", "duduk", "glissentar", "joombush", "mandolin", "mandolino",
		"wurli", "wurlitzer", "clav", "clavinet", "accordion", "chime", "chimes",
	},
	CategoryTransition: {
		"fill", "transition", "build", "buildup", "build-up", "breakdown", "break-down",
		"downlifter", "stop",
	},
	CategoryAmbiance: {
		"ambiance", "ambient", "atmosphere", "drone", "texture", "atmospheric",
	},
	CategoryFoley: {
		"foley", "bird", "animal", "water", "splash", "scratch", "vinyl", "snap",
		"whistle", "ocean", "nature", "wind",
	},
	CategoryLoop: {
		"loop", "phrase", "bar", "beat",
	},
	CategoryOneShot: {
		"oneshot", "one-shot", "hit", "stab", "shot",
	},
}

// subcategoryKeywords maps specific keywords to subcategory folder names
var subcategoryKeywords = map[Category]map[string]string{
	CategoryDrum: {
		// Existing subcategories
		"kick":  "kick",
		"bd":    "kick",
		"snare": "snare",
		"sd":    "snare",
		"hihat":       "hihat",
		"hi-hat":      "hihat",
		"hi_hat":      "hihat",
		"hi hat":      "hihat",
		"hh":          "hihat",
		"hats":        "hihat",
		"closed hat":  "hihat",
		"open hat":    "hihat",
		"hat closed":  "hihat",
		"hat open":    "hihat",
		"clap":   "clap",
		"tom":    "tom",
		"toms":   "tom",
		"cymbal": "cymbal",
		"crash":  "cymbal",
		"ride":   "cymbal",
		"cup":    "cymbal",
		"cym":    "cymbal",
		"china":  "cymbal",
		"sidestick":  "rimshot",
		"side stick": "rimshot",
		"rimshot":    "rimshot",
		"rim shot":   "rimshot",
		"crossstick": "rimshot",
		"cross stick": "rimshot",
		"rim":        "rimshot",
		// New subcategories
		"drum fill":       "fill",
		"drum_fill":       "fill",
		"drum loop":       "loop",
		"drum_loop":       "loop",
		"beat loop":       "loop",
		"beat_loop":       "loop",
		"ethnic drum":     "ethnic",
		"ethnic_drum":     "ethnic",
		"indian drum":     "ethnic",
		"indian_drum":     "ethnic",
		"tribal drum":     "ethnic",
		"tribal_drum":     "ethnic",
		"acoustic drum":   "acoustic",
		"acoustic_drum":   "acoustic",
		"cinematic drum":  "cinematic",
		"cinematic_drum":  "cinematic",
		"cinematic":       "cinematic",
	},
	CategoryBass: {
		// Existing subcategories
		"sub":      "sub",
		"subbass":  "sub",
		"sub-bass": "sub",
		"sub_bass": "sub",
		"808": "808",
		"909": "909",
		// New subcategories
		"growl":  "growl",
		"wobble": "growl",
		"whomp":  "growl",
		"freak":  "growl",
		"bass loop":  "loop",
		"bass_loop":  "loop",
		"bassloop":   "loop",
		"psy":        "psy",
		"psy bass":   "psy",
		"psy_bass":   "psy",
		"psybass":    "psy",
		"bass pluck":   "pluck",
		"bass_pluck":   "pluck",
		"pluck bass":   "pluck",
		"pluck_bass":   "pluck",
		"plucked bass": "pluck",
		"plucked_bass": "pluck",
	},
	CategorySynth: {
		// Existing subcategories
		"lead":        "lead",
		"leads":       "lead",
		"synth lead":  "lead",
		"synth_lead":  "lead",
		"pad":         "pad",
		"pads":        "pad",
		"synth pad":   "pad",
		"synth_pad":   "pad",
		"pluck":        "pluck",
		"plucks":       "pluck",
		"plucked":      "pluck",
		"synth pluck":  "pluck",
		"synth_pluck":  "pluck",
		"saw":      "saw",
		"sawtooth": "saw",
		"square": "square",
		"sine":   "sine",
		// New subcategories
		"synth loop":  "loop",
		"synth_loop":  "loop",
		"synthloop":   "loop",
		"reverse synth": "reverse",
		"reverse_synth": "reverse",
		"reversed":      "reverse",
		"synth fill":  "fill",
		"synth_fill":  "fill",
		"synthfill":   "fill",
		"arp":         "arp",
		"arpeggio":    "arp",
		"arpeggiated": "arp",
		"blip":  "blip",
		"beep":  "blip",
		"bleep": "blip",
	},
	CategoryVocal: {
		"vocal":    "vocal",
		"vox":      "vox",
		"voice":    "voice",
		"acapella": "acapella",
		"choir":    "choir",
		"chorus":   "choir",
		"ensemble": "choir",
		"shout":    "shout",
		"yell":     "shout",
		"scream":   "shout",
		"chant":    "chant",
		"chanting": "chant",
		"adlib":    "adlib",
		"ad-lib":   "adlib",
		"ad lib":   "adlib",
	},
	CategoryFX: {
		// Existing subcategories
		"riser":     "riser",
		"uplift":    "riser",
		"risefx":    "riser",
		"downsweep": "downsweep",
		"whoosh":    "whoosh",
		"impact":    "impact",
		"boom":      "impact",
		"slam":      "impact",
		"sweep":     "sweep",
		"uplifter":  "sweep",
		"noise":       "noise",
		"white":       "noise",
		"white noise": "noise",
		"pink noise":  "noise",
		"reverse": "reverse",
		"rev":     "reverse",
		// New subcategories
		"game":       "game",
		"video game": "game",
		"psy":        "psy",
		"psychedelic": "psy",
		"transformer": "transformer",
		"robot":       "transformer",
		"laser":  "laser",
		"lazer":  "laser",
		"water":  "water",
		"splash": "water",
		"ocean":  "water",
		"glitch": "glitch",
		"tone":   "tone",
		"envelope": "envelope",
		"pulse":  "pulse",
		"ufo":    "ufo",
		"bleeps": "blip",
		"sync":   "sync",
		"click":  "click",
	},
	CategoryPercussion: {
		// Existing subcategories
		"shaker": "shaker",
		"shake":  "shaker",
		"conga":  "conga",
		"congas": "conga",
		"bongo":      "bongo",
		"tambourine": "tambourine",
		"tamb":       "tambourine",
		"cowbell":    "cowbell",
		"cow bell":   "cowbell",
		"cowb":       "cowbell",
		"cabasa":     "cabasa",
		"clave":      "clave",
		"claves":     "clave",
		"agogo":      "agogo",
		"timbale":    "timbale",
		"timpani":    "timpani",
		"maracas":    "maracas",
		"maraca":     "maracas",
		"woodblock":  "woodblock",
		"wood block": "woodblock",
		"triangle":   "triangle",
		"guiro":      "guiro",
		"djembe":     "djembe",
		"udu":        "udu",
		"brush":      "brush",
		"chk":        "miscellaneous",
		// New subcategories
		"hi perc":           "high",
		"hi_perc":           "high",
		"high perc":         "high",
		"high_perc":         "high",
		"high percussion":   "high",
		"high_percussion":   "high",
		"percussion high":   "high",
		"percussion_high":   "high",
		"low perc":          "low",
		"low_perc":          "low",
		"low percussion":    "low",
		"low_percussion":    "low",
		"percussion low":    "low",
		"percussion_low":    "low",
		"mid perc":          "mid",
		"mid_perc":          "mid",
		"mid percussion":    "mid",
		"mid_percussion":    "mid",
		"percussion mid":    "mid",
		"percussion_mid":    "mid",
		"percussion loop":   "loop",
		"percussion_loop":   "loop",
		"perc loop":         "loop",
		"perc_loop":         "loop",
		"rimshot":           "rimshot",
		"rim shot":          "rimshot",
		"rim_shot":          "rimshot",
		"rim":               "rimshot",
		"clank":             "clank",
		"metal perc":        "clank",
		"metal_perc":        "clank",
		"metallic":          "clank",
		"wooden":            "wood",
		"wood perc":         "wood",
		"wood_perc":         "wood",
		"wooden perc":       "wood",
		"wooden_perc":       "wood",
		"slap":              "slap",
		"percussion slap":   "slap",
		"percussion_slap":   "slap",
		"knock":             "knock",
		"percussion knock":  "knock",
		"percussion_knock":  "knock",
		"beatbox":           "beatbox",
		"beat box":          "beatbox",
		"beat_box":          "beatbox",
		"ethnic perc":       "ethnic",
		"ethnic_perc":       "ethnic",
		"tribal perc":       "ethnic",
		"tribal_perc":       "ethnic",
		"african perc":      "ethnic",
		"african_perc":      "ethnic",
		"indian perc":       "ethnic",
		"indian_perc":       "ethnic",
	},
	CategoryMelodic: {
		// Existing subcategories
		"piano":   "piano",
		"guitar":  "guitar",
		"gtr":     "guitar",
		"acoustic guitar": "guitar",
		"electric guitar": "guitar",
		"bell":      "bell",
		"chime":     "bell",
		"chimes":    "bell",
		"marimba":   "marimba",
		"xylophone": "xylophone",
		"harp":    "harp",
		"strings": "strings",
		"string":  "strings",
		"violin":  "strings",
		"cello":   "strings",
		"viola":   "strings",
		"flute":     "woodwind",
		"clarinet":  "woodwind",
		"oboe":      "woodwind",
		"sax":       "woodwind",
		"saxophone": "woodwind",
		"horn":     "brass",
		"trumpet":  "brass",
		"trombone": "brass",
		"organ":    "keys",
		"keys":     "keys",
		"keyboard": "keys",
		"brass":    "brass",
		"woodwind": "woodwind",
		// New subcategories
		"oud":        "oud",
		"bouzouki":   "bouzouki",
		"duduk":      "duduk",
		"glissentar": "glissentar",
		"joombush":   "joombush",
		"mandolin":   "mandolin",
		"mandolino":  "mandolin",
		"wurli":      "keys",
		"wurlitzer":  "keys",
		"clav":       "keys",
		"clavinet":   "keys",
		"accordion":  "accordion",
	},
	CategoryTransition: {
		"fill":       "fill",
		"transition": "transition",
		"build":      "buildup",
		"buildup":    "buildup",
		"build-up":   "buildup",
		"breakdown":  "breakdown",
		"break-down": "breakdown",
		"downlifter": "downlifter",
		"stop":       "stop",
	},
	CategoryAmbiance: {
		"dark":        "dark",
		"bright":      "bright",
		"space":       "space",
		"nature":      "nature",
		"industrial":  "industrial",
	},
	CategoryFoley: {
		"bird":      "nature",
		"animal":    "animal",
		"water":     "water",
		"splash":    "water",
		"ocean":     "water",
		"scratch":   "vinyl",
		"vinyl":     "vinyl",
		"snap":      "human",
		"whistle":   "human",
		"wind":      "nature",
		"mechanical": "mechanical",
	},
	CategoryLoop: {
		"loop":   "loop",
		"phrase": "phrase",
		"bar":    "bar",
		"beat":   "beat",
	},
	CategoryOneShot: {
		// Restructured for instrument-specific subcategories
		"bass shot":    "bass",
		"bass_shot":    "bass",
		"bass stab":    "bass",
		"bass_stab":    "bass",
		"bass hit":     "bass",
		"bass_hit":     "bass",
		"bassshot":     "bass",
		"synth shot":   "synth",
		"synth_shot":   "synth",
		"synth stab":   "synth",
		"synth_stab":   "synth",
		"synthshot":    "synth",
		"vocal shot":   "vocal",
		"vocal_shot":   "vocal",
		"drum hit":     "drum",
		"drum_hit":     "drum",
		"drum stab":    "drum",
		"drum_stab":    "drum",
		"melodic stab": "melodic",
		"melodic_stab": "melodic",
		"oneshot":      "general",
		"one-shot":     "general",
		"one_shot":     "general",
		"hit":          "general",
		"stab":         "general",
		"shot":         "general",
	},
}

// NormalizeFileName normalizes a filename by applying transformations:
// - converts to lowercase
// - replaces spaces with dashes
// - replaces underscores with dashes
// - removes consecutive dashes
// - preserves the file extension
func NormalizeFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	nameWithoutExt := strings.TrimSuffix(fileName, ext)

	// Apply normalization transformations
	normalized := strings.ToLower(nameWithoutExt)           // lowercase
	normalized = strings.ReplaceAll(normalized, " ", "-")   // spaces to dashes
	normalized = strings.ReplaceAll(normalized, "_", "-")   // underscores to dashes

	// Remove consecutive dashes
	dashRegex := regexp.MustCompile(`-+`)
	normalized = dashRegex.ReplaceAllString(normalized, "-")

	// Remove leading/trailing dashes
	normalized = strings.Trim(normalized, "-")

	return normalized + ext
}

// determineSubcategory checks the filename for subcategory keywords
// Returns the subcategory based on the longest matching keyword (most specific match)
func determineSubcategory(category Category, fileName, nameWithoutExt string) string {
	// Get subcategory map for this category
	subcatMap, exists := subcategoryKeywords[category]
	if !exists {
		return ""
	}

	// Find all matching keywords and prefer the longest one (most specific)
	var bestMatch string
	var bestKeyword string
	for keyword, subfolder := range subcatMap {
		if strings.Contains(fileName, keyword) || strings.Contains(nameWithoutExt, keyword) {
			// Prefer longer keywords (more specific matches)
			if len(keyword) > len(bestKeyword) {
				bestKeyword = keyword
				bestMatch = subfolder
			}
		}
	}

	return bestMatch
}

// Categorize determines the category of a sample file based on its name
func Categorize(sample scanner.SampleFile, targetDir string, normalize bool) CategorizedFile {
	fileName := strings.ToLower(sample.FileName)
	nameWithoutExt := strings.ToLower(strings.TrimSuffix(sample.FileName, sample.Extension))

	category := determineCategory(fileName, nameWithoutExt)
	subcategory := determineSubcategory(category, fileName, nameWithoutExt)

	// If no subcategory found but category has subcategory support, use "uncategorized"
	if subcategory == "" && category != CategoryUncategorized {
		_, hasSubcategories := subcategoryKeywords[category]
		if hasSubcategories {
			subcategory = "uncategorized"
		}
	}

	// Determine the target filename (with optional normalization)
	targetFileName := sample.FileName
	if normalize {
		targetFileName = NormalizeFileName(sample.FileName)
	}

	// Build target path with subcategory
	var targetPath string
	if subcategory != "" {
		targetPath = filepath.Join(targetDir, string(category), subcategory, targetFileName)
	} else {
		targetPath = filepath.Join(targetDir, string(category), targetFileName)
	}

	return CategorizedFile{
		Sample:      sample,
		Category:    category,
		Subcategory: subcategory,
		TargetPath:  targetPath,
	}
}

// determineCategory checks the filename against category keywords in priority order
func determineCategory(fileName, nameWithoutExt string) Category {
	// Check categories in priority order to ensure deterministic results
	for _, cat := range categoryPriority {
		words := keywords[cat]
		for _, word := range words {
			if strings.Contains(fileName, word) || strings.Contains(nameWithoutExt, word) {
				return cat
			}
		}
	}
	return CategoryUncategorized
}

// CategorizeBatch categorizes multiple sample files
func CategorizeBatch(samples []scanner.SampleFile, targetDir string, normalize bool) []CategorizedFile {
	categorized := make([]CategorizedFile, 0, len(samples))

	for _, sample := range samples {
		categorized = append(categorized, Categorize(sample, targetDir, normalize))
	}

	return categorized
}
