package audioconverter

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WhisperCompatibleWAVParams represents the audio parameters known to work well with whisper.cpp
type WhisperCompatibleWAVParams struct {
	SampleRate    int    // 16000 Hz is recommended for whisper.cpp
	Channels      int    // 1 channel (mono) is recommended for whisper.cpp
	AudioFormat   string // PCM format is recommended
	BitDepth      int    // 16-bit depth is recommended
	TempDirectory string // Directory for temporary files
}

// DefaultWhisperParams returns the default parameters optimized for whisper.cpp
func DefaultWhisperParams() WhisperCompatibleWAVParams {
	return WhisperCompatibleWAVParams{
		SampleRate:    16000,
		Channels:      1,
		AudioFormat:   "pcm_s16le", // Signed 16-bit little-endian PCM
		BitDepth:      16,
		TempDirectory: os.TempDir(),
	}
}

// ConvertToWAV converts any audio file to a WAV file format compatible with whisper.cpp
// It returns the path to the converted WAV file and any error encountered
func ConvertToWAV(inputPath string, params WhisperCompatibleWAVParams) (string, error) {
	// Check if the input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("input file does not exist: %s", inputPath)
	}

	// Check if ffmpeg is available
	if err := checkFFmpeg(); err != nil {
		return "", err
	}

	// Create output file path
	outputPath := filepath.Join(params.TempDirectory,
		fmt.Sprintf("%s_whisper_compatible.wav", strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))))

	// Build the ffmpeg command
	args := []string{
		"-i", inputPath,
		"-ar", fmt.Sprintf("%d", params.SampleRate),
		"-ac", fmt.Sprintf("%d", params.Channels),
		"-c:a", params.AudioFormat,
		"-b:a", fmt.Sprintf("%dk", params.BitDepth*params.SampleRate*params.Channels/1000),
		outputPath,
		"-y", // Overwrite output file if it exists
	}

	// Execute ffmpeg
	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg conversion failed: %v, output: %s", err, string(output))
	}

	// Verify the output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return "", errors.New("conversion completed but output file was not created")
	}

	return outputPath, nil
}

// ConvertToWAVWithDefaultParams converts any audio file to a WAV file with default whisper.cpp parameters
func ConvertToWAVWithDefaultParams(inputPath string) (string, error) {
	return ConvertToWAV(inputPath, DefaultWhisperParams())
}

// checkFFmpeg verifies that ffmpeg is installed and available
func checkFFmpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		return errors.New("ffmpeg is not installed or not found in PATH")
	}
	return nil
}

// CleanupTempFile removes a temporary file created during conversion
func CleanupTempFile(filePath string) error {
	return os.Remove(filePath)
}
