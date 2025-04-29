package audioconverter

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestConvertToWAV is an example test function to exercise the audio conversion functionality
func TestConvertToWAV(t *testing.T) {
	// Create a temporary test file if no real audio file is available
	tempTestFile := filepath.Join(os.TempDir(), "test_audio.mp3")
	// Skip this part if you want to use a real audio file instead
	if _, err := os.Stat(tempTestFile); os.IsNotExist(err) {
		// Just create an empty file for testing the function execution path
		f, err := os.Create(tempTestFile)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		f.Close()
		defer os.Remove(tempTestFile) // Clean up after test
		t.Logf("Created temporary test file: %s", tempTestFile)
	}

	// Hardcoded test audio file path
	testAudioFile := "./files/xx.waptt.opus"
	t.Logf("Using specified test file: %s", testAudioFile)

	// Get the directory of the input file to use as output directory
	inputDir := filepath.Dir(testAudioFile)
	t.Logf("Using input directory for output: %s", inputDir)

	// Test with custom parameters to output in the same directory as input
	t.Run("OutputInSameDirectory", func(t *testing.T) {
		params := WhisperCompatibleWAVParams{
			SampleRate:    16000,
			Channels:      1,
			AudioFormat:   "pcm_s16le",
			BitDepth:      16,
			TempDirectory: inputDir, // Set output directory to input directory
		}

		outputFile, err := ConvertToWAV(testAudioFile, params)
		if err != nil {
			t.Fatalf("Conversion error: %v", err)
		}

		t.Logf("Successfully converted file to: %s", outputFile)

		// Verify the output file exists
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Errorf("Output file was not created: %s", outputFile)
		} else {
			t.Logf("Output file exists in same directory as input: %s", outputFile)
			// Note: Not cleaning up the file since we want to keep it
			// If you want to clean it up, uncomment the following:
			// if err := CleanupTempFile(outputFile); err != nil {
			//     t.Logf("Failed to cleanup output file: %v", err)
			// }
		}
	})

	// Test with default parameters (this will output to system temp directory)
	t.Run("DefaultParameters", func(t *testing.T) {
		outputFile, err := ConvertToWAVWithDefaultParams(testAudioFile)
		if err != nil {
			t.Fatalf("Conversion error: %v", err)
		}

		t.Logf("Successfully converted file to: %s", outputFile)

		// Verify the output file exists
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Errorf("Output file was not created: %s", outputFile)
		} else {
			t.Logf("Output file exists: %s", outputFile)
			// Clean up the file in temp directory
			if err := CleanupTempFile(outputFile); err != nil {
				t.Logf("Failed to cleanup temp file: %v", err)
			}
		}
	})

	// Test error handling with non-existent file
	t.Run("NonExistentFile", func(t *testing.T) {
		nonExistentFile := filepath.Join(os.TempDir(), "this_file_does_not_exist.mp3")
		_, err := ConvertToWAVWithDefaultParams(nonExistentFile)
		if err == nil {
			t.Errorf("Expected error for non-existent file, but got nil")
		} else {
			t.Logf("Correctly received error for non-existent file: %v", err)
		}
	})

	// Test FFmpeg check
	t.Run("FFmpegCheck", func(t *testing.T) {
		err := checkFFmpeg()
		if err != nil {
			t.Logf("FFmpeg check failed: %v", err)
		} else {
			t.Logf("FFmpeg is available in the system")
		}
	})

	fmt.Println("Test exercises completed")
}
