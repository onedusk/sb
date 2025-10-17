package mov_to_mp4

// MP4Options contains MP4-specific conversion options
type MP4Options struct {
	// Quality
	CRF    int    // 0-51, lower = better quality (default: 23)
	Preset string // ultrafast, fast, medium, slow, veryslow

	// Codecs
	VideoCodec   string // h264, h265, vp9
	AudioCodec   string // aac, mp3, copy
	AudioBitrate string // e.g., "128k", "192k"

	// Hardware
	HWAccel       string // videotoolbox, nvenc, qsv
	HWAccelDevice string

	// Bitrate control
	VideoBitrate string // e.g., "2M", "5M"
}

// DefaultMP4Options returns default options for MP4 conversion
func DefaultMP4Options() MP4Options {
	return MP4Options{
		CRF:          23,
		Preset:       "medium",
		VideoCodec:   "h264",
		AudioCodec:   "aac",
		AudioBitrate: "192k",
	}
}

// Validate checks if options are valid
func (o *MP4Options) Validate() error {
	// CRF validation
	if o.CRF < 0 || o.CRF > 51 {
		o.CRF = 23
	}

	// Preset validation
	validPresets := map[string]bool{
		"ultrafast": true,
		"superfast": true,
		"veryfast":  true,
		"faster":    true,
		"fast":      true,
		"medium":    true,
		"slow":      true,
		"slower":    true,
		"veryslow":  true,
	}

	if o.Preset != "" && !validPresets[o.Preset] {
		o.Preset = "medium"
	}

	return nil
}
