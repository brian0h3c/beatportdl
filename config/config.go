package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"unspok3n/beatportdl/internal/validator"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Username      string `yaml:"username,omitempty"`
	Password      string `yaml:"password,omitempty"`
	Quality       string `yaml:"quality,omitempty"`
	WriteErrorLog bool   `yaml:"write_error_log,omitempty"`
	ShowProgress  bool   `yaml:"show_progress,omitempty"`

	MaxGlobalWorkers   int `yaml:"max_global_workers,omitempty"`
	MaxDownloadWorkers int `yaml:"max_download_workers,omitempty"`

	DownloadsDirectory      string `yaml:"downloads_directory,omitempty"`
	SortByContext           bool   `yaml:"sort_by_context,omitempty"`
	SortByLabel             bool   `yaml:"sort_by_label,omitempty"`
	ForceReleaseDirectories bool   `yaml:"force_release_directories,omitempty"`
	TrackExists             string `yaml:"track_exists,omitempty"`
	TrackNumberPadding      int    `yaml:"track_number_padding,omitempty"`

	ReleaseDirectoryTemplate  string `yaml:"release_directory_template,omitempty"`
	PlaylistDirectoryTemplate string `yaml:"playlist_directory_template,omitempty"`
	ChartDirectoryTemplate    string `yaml:"chart_directory_template,omitempty"`
	LabelDirectoryTemplate    string `yaml:"label_directory_template,omitempty"`
	ArtistDirectoryTemplate   string `yaml:"artist_directory_template,omitempty"`
	TrackFileTemplate         string `yaml:"track_file_template,omitempty"`
	WhitespaceCharacter       string `yaml:"whitespace_character,omitempty"`
	ArtistsLimit              int    `yaml:"artists_limit,omitempty"`
	ArtistsShortForm          string `yaml:"artists_short_form,omitempty"`
	KeySystem                 string `yaml:"key_system,omitempty"`

	CoverSize string `yaml:"cover_size,omitempty"`
	KeepCover bool   `yaml:"keep_cover,omitempty"`
	FixTags   bool   `yaml:"fix_tags,omitempty"`

	TagMappings map[string]map[string]string `yaml:"tag_mappings,omitempty"`

	Proxy string `yaml:"proxy,omitempty"`
}

const (
	DefaultCoverSize = "1400x1400"
)

// qualityCompatibilityComment is written to the top of a generated config file
// to document the quality values and which Pioneer DJ / AlphaTheta hardware can
// play each one.
const qualityCompatibilityComment = `# quality options (Pioneer DJ / AlphaTheta player compatibility):
#   lossless    44.1 kHz FLAC   rekordbox + modern players (CDJ-3000, CDJ-2000NXS2, XDJ-XZ/RX3/RR/AZ, OPUS-QUAD); not older CDJs
#   high        256 kbps AAC    rekordbox + virtually all Pioneer DJ / AlphaTheta gear (incl. older CDJs)
#   medium      128 kbps AAC    same broad compatibility as high, lower bitrate
#   medium-hls  128 kbps AAC    same as medium (requires ffmpeg)
#
# downloads_directory defaults to a "downloads" folder in the current directory
# (where beatportdl is run) when left empty.

`

var (
	SupportedTrackExistsOptions = []string{
		"error",
		"skip",
		"overwrite",
		"update",
	}

	SupportedKeySystems = []string{
		"standard",
		"standard-short",
		"openkey",
		"camelot",
	}
)

func FFMPEGInstalled() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

func Parse(filePath string) (*AppConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	config := AppConfig{
		Quality:                   "lossless",
		CoverSize:                 DefaultCoverSize,
		TrackFileTemplate:         "{number}. {artists} - {name} ({mix_name})",
		ReleaseDirectoryTemplate:  "[{catalog_number}] {artists} - {name}",
		PlaylistDirectoryTemplate: "{name}",
		ChartDirectoryTemplate:    "{name} [{published_date}]",
		LabelDirectoryTemplate:    "{name} [{updated_date}]",
		ArtistDirectoryTemplate:   "{name}",
		SortByContext:             true,
		ArtistsLimit:              3,
		ArtistsShortForm:          "VA",
		KeySystem:                 "standard-short",
		TrackExists:               "update",
		TrackNumberPadding:        2,
		FixTags:                   true,
		ShowProgress:              true,
		MaxGlobalWorkers:          15,
		MaxDownloadWorkers:        15,
	}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	if config.Username == "" || config.Password == "" {
		return nil, fmt.Errorf("username or password is not provided")
	}

	if config.Quality == "medium-hls" && !FFMPEGInstalled() {
		return nil, errors.New("ffmpeg not found")
	}

	if config.TagMappings != nil {
		if err = ValidateTagMappings(config.TagMappings); err != nil {
			return nil, err
		}

		if _, ok := config.TagMappings["flac"]; !ok {
			config.TagMappings["flac"] = DefaultTagMappings["flac"]
		}

		if _, ok := config.TagMappings["m4a"]; !ok {
			config.TagMappings["m4a"] = DefaultTagMappings["m4a"]
		}
	} else {
		config.TagMappings = DefaultTagMappings
	}

	if !validator.PermittedValue(config.KeySystem, SupportedKeySystems...) {
		return nil, fmt.Errorf("invalid key system")
	}

	if config.DownloadsDirectory == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("no downloads directory provided and failed to resolve current directory: %w", err)
		}
		config.DownloadsDirectory = filepath.Join(cwd, "downloads")
	}

	if !validator.PermittedValue(config.TrackExists, SupportedTrackExistsOptions...) {
		return nil, fmt.Errorf("invalid track exists behavior")
	}

	if config.TrackNumberPadding > 10 || config.TrackNumberPadding < 0 {
		return nil, fmt.Errorf("invalid track number padding")
	}

	return &config, nil
}

func (c *AppConfig) Save(filePath string) error {
	err := os.MkdirAll(path.Dir(filePath), 0700)
	if err != nil {
		return fmt.Errorf("could not create folder for configuration file: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(qualityCompatibilityComment); err != nil {
		return fmt.Errorf("write config header: %w", err)
	}

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(&c); err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	return nil
}
