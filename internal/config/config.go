package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// CategoryConfig represents the configuration for categories and subcategories
type CategoryConfig struct {
	Categories []CategoryDefinition `json:"categories"`
}

// CategoryDefinition defines a single category with its keywords and subcategories
type CategoryDefinition struct {
	Name         string                `json:"name"`
	Priority     int                   `json:"priority"`
	Keywords     []string              `json:"keywords"`
	Subcategories map[string][]string  `json:"subcategories,omitempty"`
}

// LoadConfig loads the configuration from a JSON file
// If configPath is empty, returns the default configuration
func LoadConfig(configPath string) (*CategoryConfig, error) {
	if configPath == "" {
		return GetDefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config CategoryConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// validateConfig checks that the given CategoryConfig is valid.
// It ensures that there is at least one category, each category has a non-empty name,
// no duplicate category names exist, and each category has at least one keyword.
// Returns an error describing the first validation failure encountered, or nil if valid.
func validateConfig(config *CategoryConfig) error {
	if len(config.Categories) == 0 {
		return fmt.Errorf("configuration must contain at least one category")
	}

	seenNames := make(map[string]bool)
	for _, cat := range config.Categories {
		if cat.Name == "" {
			return fmt.Errorf("category name cannot be empty")
		}
		if seenNames[cat.Name] {
			return fmt.Errorf("duplicate category name: %s", cat.Name)
		}
		seenNames[cat.Name] = true

		if len(cat.Keywords) == 0 {
			return fmt.Errorf("category %s must have at least one keyword", cat.Name)
		}
	}

	return nil
}

// GetDefaultConfig returns the default hardcoded configuration
// This ensures backward compatibility when no config file is provided
func GetDefaultConfig() *CategoryConfig {
	return &CategoryConfig{
		Categories: []CategoryDefinition{
			{
				Name:     "oneshots",
				Priority: 1,
				Keywords: []string{"oneshot", "one-shot", "hit", "stab", "shot"},
				Subcategories: map[string][]string{
					"bass":    {"bass shot", "bass_shot", "bass stab", "bass_stab", "bass hit", "bass_hit", "bassshot"},
					"synth":   {"synth shot", "synth_shot", "synth stab", "synth_stab", "synthshot"},
					"vocal":   {"vocal shot", "vocal_shot"},
					"drum":    {"drum hit", "drum_hit", "drum stab", "drum_stab"},
					"melodic": {"melodic stab", "melodic_stab"},
					"general": {"oneshot", "one-shot", "one_shot", "hit", "stab", "shot"},
				},
			},
			{
				Name:     "drums",
				Priority: 2,
				Keywords: []string{"kick", "snare", "hihat", "hi-hat", "hi_hat", "hi hat", "hats", "clap", "tom", "cymbal", "crash", "ride", "drum", "bd", "sd", "hh", "closed hat", "open hat", "hat closed", "hat open", "sidestick", "side stick", "rimshot", "rim shot", "cup", "rim", "cym", "china", "crossstick", "cross stick"},
				Subcategories: map[string][]string{
					"kick":   {"kick", "bd"},
					"snare":  {"snare", "sd"},
					"hihat":  {"hihat", "hi-hat", "hi_hat", "hi hat", "hh", "hats", "closed hat", "open hat", "hat closed", "hat open"},
					"clap":   {"clap"},
					"tom":    {"tom", "toms"},
					"cymbal": {"cymbal", "crash", "ride", "cup", "cym", "china"},
					"rimshot": {"sidestick", "side stick", "rimshot", "rim shot", "crossstick", "cross stick", "rim"},
					"fill":   {"drum fill", "drum_fill"},
					"loop":   {"drum loop", "drum_loop", "beat loop", "beat_loop"},
					"ethnic": {"ethnic drum", "ethnic_drum", "indian drum", "indian_drum", "tribal drum", "tribal_drum"},
					"acoustic": {"acoustic drum", "acoustic_drum"},
					"cinematic": {"cinematic drum", "cinematic_drum", "cinematic"},
				},
			},
			{
				Name:     "bass",
				Priority: 3,
				Keywords: []string{"bass", "sub", "808", "909"},
				Subcategories: map[string][]string{
					"sub":   {"sub", "subbass", "sub-bass", "sub_bass"},
					"808":   {"808"},
					"909":   {"909"},
					"growl": {"growl", "wobble", "whomp", "freak"},
					"loop":  {"bass loop", "bass_loop", "bassloop"},
					"psy":   {"psy", "psy bass", "psy_bass", "psybass"},
					"pluck": {"bass pluck", "bass_pluck", "pluck bass", "pluck_bass", "plucked bass", "plucked_bass"},
				},
			},
			{
				Name:     "percussion",
				Priority: 4,
				Keywords: []string{"perc", "percussion", "shaker", "conga", "bongo", "tambourine", "tamb", "cowbell", "cabasa", "clave", "claves", "agogo", "timbale", "timpani", "maracas", "maraca", "woodblock", "wood block", "triangle", "guiro", "djembe", "udu", "brush", "chk", "cowb"},
				Subcategories: map[string][]string{
					"shaker":      {"shaker", "shake"},
					"conga":       {"conga", "congas"},
					"bongo":       {"bongo"},
					"tambourine":  {"tambourine", "tamb"},
					"cowbell":     {"cowbell", "cow bell", "cowb"},
					"cabasa":      {"cabasa"},
					"clave":       {"clave", "claves"},
					"agogo":       {"agogo"},
					"timbale":     {"timbale"},
					"timpani":     {"timpani"},
					"maracas":     {"maracas", "maraca"},
					"woodblock":   {"woodblock", "wood block"},
					"triangle":    {"triangle"},
					"guiro":       {"guiro"},
					"djembe":      {"djembe"},
					"udu":         {"udu"},
					"brush":       {"brush"},
					"miscellaneous": {"chk"},
					"high":        {"hi perc", "hi_perc", "high perc", "high_perc", "high percussion", "high_percussion", "percussion high", "percussion_high"},
					"low":         {"low perc", "low_perc", "low percussion", "low_percussion", "percussion low", "percussion_low"},
					"mid":         {"mid perc", "mid_perc", "mid percussion", "mid_percussion", "percussion mid", "percussion_mid"},
					"loop":        {"percussion loop", "percussion_loop", "perc loop", "perc_loop"},
					"rimshot":     {"rimshot", "rim shot", "rim_shot", "rim"},
					"clank":       {"clank", "metal perc", "metal_perc", "metallic"},
					"wood":        {"wooden", "wood perc", "wood_perc", "wooden perc", "wooden_perc"},
					"slap":        {"slap", "percussion slap", "percussion_slap"},
					"knock":       {"knock", "percussion knock", "percussion_knock"},
					"beatbox":     {"beatbox", "beat box", "beat_box"},
					"ethnic":      {"ethnic perc", "ethnic_perc", "tribal perc", "tribal_perc", "african perc", "african_perc", "indian perc", "indian_perc"},
				},
			},
			{
				Name:     "vocals",
				Priority: 5,
				Keywords: []string{"vocal", "vox", "voice", "acapella", "choir", "shout", "chant", "adlib"},
				Subcategories: map[string][]string{
					"vocal":    {"vocal"},
					"vox":      {"vox"},
					"voice":    {"voice"},
					"acapella": {"acapella"},
					"choir":    {"choir", "chorus", "ensemble"},
					"shout":    {"shout", "yell", "scream"},
					"chant":    {"chant", "chanting"},
					"adlib":    {"adlib", "ad-lib", "ad lib"},
				},
			},
			{
				Name:     "synth",
				Priority: 6,
				Keywords: []string{"synth", "lead", "pad", "pluck", "saw", "square", "sine"},
				Subcategories: map[string][]string{
					"lead":    {"lead", "leads", "synth lead", "synth_lead"},
					"pad":     {"pad", "pads", "synth pad", "synth_pad"},
					"pluck":   {"pluck", "plucks", "plucked", "synth pluck", "synth_pluck"},
					"saw":     {"saw", "sawtooth"},
					"square":  {"square"},
					"sine":    {"sine"},
					"loop":    {"synth loop", "synth_loop", "synthloop"},
					"reverse": {"reverse synth", "reverse_synth", "reversed"},
					"fill":    {"synth fill", "synth_fill", "synthfill"},
					"arp":     {"arp", "arpeggio", "arpeggiated"},
					"blip":    {"blip", "beep", "bleep"},
				},
			},
			{
				Name:     "melodic",
				Priority: 7,
				Keywords: []string{"piano", "guitar", "bell", "marimba", "xylophone", "harp", "strings", "violin", "cello", "flute", "horn", "trumpet", "sax", "saxophone", "organ", "keys", "brass", "woodwind", "arpeggio", "arpeggiated", "melody", "oud", "bouzouki", "duduk", "glissentar", "joombush", "mandolin", "mandolino", "wurli", "wurlitzer", "clav", "clavinet", "accordion", "chime", "chimes"},
				Subcategories: map[string][]string{
					"piano":     {"piano"},
					"guitar":    {"guitar", "gtr", "acoustic guitar", "electric guitar"},
					"bell":      {"bell", "chime", "chimes"},
					"marimba":   {"marimba"},
					"xylophone": {"xylophone"},
					"harp":      {"harp"},
					"strings":   {"strings", "string", "violin", "cello", "viola"},
					"woodwind":  {"flute", "clarinet", "oboe", "sax", "saxophone", "woodwind"},
					"brass":     {"horn", "trumpet", "trombone", "brass"},
					"keys":      {"organ", "keys", "keyboard", "wurli", "wurlitzer", "clav", "clavinet"},
					"oud":       {"oud"},
					"bouzouki":  {"bouzouki"},
					"duduk":     {"duduk"},
					"glissentar": {"glissentar"},
					"joombush":  {"joombush"},
					"mandolin":  {"mandolin", "mandolino"},
					"accordion": {"accordion"},
				},
			},
			{
				Name:     "fx",
				Priority: 8,
				Keywords: []string{"fx", "sfx", "riser", "downsweep", "whoosh", "impact", "sweep", "noise", "white", "reverse", "rev", "glitch", "tone", "envelope", "pulse", "ufo", "bleeps", "sync", "click"},
				Subcategories: map[string][]string{
					"riser":       {"riser", "uplift", "risefx"},
					"downsweep":   {"downsweep"},
					"whoosh":      {"whoosh"},
					"impact":      {"impact", "boom", "slam"},
					"sweep":       {"sweep", "uplifter"},
					"noise":       {"noise", "white", "white noise", "pink noise"},
					"reverse":     {"reverse", "rev"},
					"game":        {"game", "video game"},
					"psy":         {"psy", "psychedelic"},
					"transformer": {"transformer", "robot"},
					"laser":       {"laser", "lazer"},
					"water":       {"water", "splash", "ocean"},
					"glitch":      {"glitch"},
					"tone":        {"tone"},
					"envelope":    {"envelope"},
					"pulse":       {"pulse"},
					"ufo":         {"ufo"},
					"blip":        {"bleeps"},
					"sync":        {"sync"},
					"click":       {"click"},
				},
			},
			{
				Name:     "transition",
				Priority: 9,
				Keywords: []string{"fill", "transition", "build", "buildup", "build-up", "breakdown", "break-down", "downlifter", "stop"},
				Subcategories: map[string][]string{
					"fill":       {"fill"},
					"transition": {"transition"},
					"buildup":    {"build", "buildup", "build-up"},
					"breakdown":  {"breakdown", "break-down"},
					"downlifter": {"downlifter"},
					"stop":       {"stop"},
				},
			},
			{
				Name:     "ambiance",
				Priority: 10,
				Keywords: []string{"ambiance", "ambient", "atmosphere", "drone", "texture", "atmospheric"},
				Subcategories: map[string][]string{
					"dark":       {"dark"},
					"bright":     {"bright"},
					"space":      {"space"},
					"nature":     {"nature"},
					"industrial": {"industrial"},
				},
			},
			{
				Name:     "foley",
				Priority: 11,
				Keywords: []string{"foley", "bird", "animal", "water", "splash", "scratch", "vinyl", "snap", "whistle", "ocean", "nature", "wind"},
				Subcategories: map[string][]string{
					"nature":     {"bird", "wind"},
					"animal":     {"animal"},
					"water":      {"water", "splash", "ocean"},
					"vinyl":      {"scratch", "vinyl"},
					"human":      {"snap", "whistle"},
					"mechanical": {"mechanical"},
				},
			},
			{
				Name:     "loops",
				Priority: 12,
				Keywords: []string{"loop", "phrase", "bar", "beat"},
				Subcategories: map[string][]string{
					"loop":   {"loop"},
					"phrase": {"phrase"},
					"bar":    {"bar"},
					"beat":   {"beat"},
				},
			},
		},
	}
}
