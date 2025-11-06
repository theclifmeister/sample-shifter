# Sample Shifter

A simple CLI tool to organize audio sample files by automatically categorizing them based on their names.

## Features

- **Automatic Categorization**: Intelligently categorizes samples based on filename keywords
- **Non-Destructive**: Original files remain untouched; copies are made to the target directory
- **Preview Mode**: See how files will be organized before applying changes
- **Multiple Audio Formats**: Supports WAV, MP3, FLAC, AIF, AIFF, OGG, M4A, WMA, and AAC
- **Recursive Scanning**: Finds samples in subdirectories
- **Flexible Workflow**: Preview, save, and apply later or apply directly

## Categories

Sample Shifter automatically categorizes files into the following categories:

- **drums**: kicks, snares, hi-hats, claps, toms, cymbals, etc.
- **bass**: bass sounds, sub bass, 808s, 909s
- **synth**: synth leads, pads, plucks, and other synthesizer sounds
- **vocals**: vocal shots, acapellas, choirs, shouts
- **fx**: risers, downsweps, whooshes, impacts, noise
- **percussion**: shakers, congas, bongos, tambourines, cowbells
- **loops**: drum loops, phrases, beats
- **oneshots**: one-shot samples, hits, stabs
- **uncategorized**: files that don't match any category

## Installation

### Prerequisites

- Go 1.24 or later

### Build from Source

```bash
git clone https://github.com/theclifmeister/sample-shifter.git
cd sample-shifter
go build -o sample-shifter
```

## Usage

### Basic Commands

#### 1. Scan Directory

Scan a directory to see what audio files are present:

```bash
./sample-shifter scan /path/to/samples
```

#### 2. Preview Categorization

Preview how files will be organized without making any changes:

```bash
./sample-shifter preview /path/to/samples --target /path/to/organized
```

Save the preview to a file for later use:

```bash
./sample-shifter preview /path/to/samples --target /path/to/organized --output preview.json
```

#### 3. Apply Changes

Apply the categorization and copy files to the target directory:

```bash
./sample-shifter apply /path/to/samples --target /path/to/organized
```

Use a previously saved preview file:

```bash
./sample-shifter apply --preview-file preview.json
```

Perform a dry run to see what would happen without actually copying files:

```bash
./sample-shifter apply /path/to/samples --target /path/to/organized --dry-run
```

### Typical Workflow

1. **Scan** your sample library to see what files are present:
   ```bash
   ./sample-shifter scan ~/Music/Samples
   ```

2. **Preview** the categorization:
   ```bash
   ./sample-shifter preview ~/Music/Samples --target ~/Music/Organized --output preview.json
   ```

3. Review the preview output to ensure categorization is correct

4. **Apply** the changes:
   ```bash
   ./sample-shifter apply --preview-file preview.json
   ```

### Command Reference

#### `scan [directory]`

Scans a directory recursively for audio sample files.

**Arguments:**
- `directory`: Path to the directory to scan

**Example:**
```bash
./sample-shifter scan /path/to/samples
```

#### `preview [source-directory]`

Previews how files will be categorized and organized.

**Arguments:**
- `source-directory`: Path to the source directory containing samples

**Flags:**
- `--target, -t`: Target directory for organized samples (required)
- `--output, -o`: Save preview to JSON file

**Example:**
```bash
./sample-shifter preview /path/to/samples --target /path/to/organized --output preview.json
```

#### `apply [source-directory]`

Applies categorization and copies files to the target directory.

**Arguments:**
- `source-directory`: Path to the source directory (optional if using --preview-file)

**Flags:**
- `--target, -t`: Target directory for organized samples (required if not using --preview-file)
- `--preview-file, -p`: Use a previously saved preview file
- `--dry-run`: Preview what would be done without actually copying files

**Examples:**
```bash
# Apply directly from source
./sample-shifter apply /path/to/samples --target /path/to/organized

# Apply from preview file
./sample-shifter apply --preview-file preview.json

# Dry run
./sample-shifter apply /path/to/samples --target /path/to/organized --dry-run
```

## Examples

### Organize a Sample Library

```bash
# Preview the organization
./sample-shifter preview ~/Downloads/NewSamples --target ~/Music/Organized

# Apply after reviewing
./sample-shifter apply ~/Downloads/NewSamples --target ~/Music/Organized
```

### Using Preview Files for Review

```bash
# Create a preview file
./sample-shifter preview ~/Downloads/NewSamples --target ~/Music/Organized --output review.json

# Review the JSON file, make adjustments if needed

# Apply the categorization
./sample-shifter apply --preview-file review.json
```

## How Categorization Works

Sample Shifter uses keyword matching on filenames to determine categories. For example:

- `kick_01.wav` → **drums** (contains "kick")
- `bass_sub_heavy.wav` → **bass** (contains "bass" and "sub")
- `synth_lead_bright.flac` → **synth** (contains "synth" and "lead")
- `vocal_shot.mp3` → **vocals** (contains "vocal")
- `fx_riser.wav` → **fx** (contains "fx" and "riser")

Files that don't match any keywords are placed in the **uncategorized** folder.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the MIT License.
