package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestImgToRainData(t *testing.T) {
	// Load the test image.
	// This assumes you have a 'test_files/radar_rain.png' file for testing.
	testImagePath := filepath.Join("test_files", "radar_rain.png")
	data, err := os.ReadFile(testImagePath)
	if err != nil {
		t.Fatalf("Failed to read test image '%s': %v. Please ensure it exists.", testImagePath, err)
	}

	// Call the function to be tested.
	rainData := imgToRainData(data)

	// --- Assertions ---

	// 1. The function should return non-nil data for a valid image.
	if rainData == nil {
		t.Fatal("imgToRainData returned nil for a valid image.")
	}

	// 2. The dimensions of the rain data should match the image dimensions.
	// The DMI images are typically 512x512.
	expectedWidth := 512
	expectedHeight := 512
	if len(rainData) != expectedHeight {
		t.Errorf("Expected rainData height to be %d, but got %d", expectedHeight, len(rainData))
	}
	if len(rainData) > 0 && len(rainData[0]) != expectedWidth {
		t.Errorf("Expected rainData width to be %d, but got %d", expectedWidth, len(rainData[0]))
	}

	// 3. Check specific pixel values.
	testCases := []struct {
		name          string
		x, y          int
		expectedIndex int // Index from the 'rainColors' slice
	}{
		// Example: A pixel with heavy rain. Color {128, 0, 0} is at index 1.
		{"heavy rain", 264, 214, 8},
		// Example: A pixel with light rain. Color {158, 242, 233} is at index 15.
		{"light rain", 249, 213, 2},
		// Example: A pixel with no rain (a color not in rainColors should result in 0).
		{"no rain", 250, 210, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.y >= len(rainData) || tc.x >= len(rainData[tc.y]) {
				t.Fatalf("Test coordinate (%d, %d) is out of bounds for the image.", tc.x, tc.y)
			}
			if got := rainData[tc.y][tc.x]; got != tc.expectedIndex {
				t.Errorf("At pixel (%d, %d), expected rain index %d, but got %d", tc.x, tc.y, tc.expectedIndex, got)
			}
		})
	}
}

// TestPrintRainDataForVisualInspection processes the test image and prints a
// subsection of the resulting rain data grid to the console. This is useful for
// manually verifying that the color-to-index mapping is working as expected.
// To see the output, run tests with the verbose flag: `go test -v`.
func TestPrintRainDataForVisualInspection(t *testing.T) {
	testImagePath := filepath.Join("test_files", "radar_rain.png")
	data, err := os.ReadFile(testImagePath)
	if err != nil {
		t.Fatalf("Failed to read test image '%s': %v", testImagePath, err)
	}

	rainData := imgToRainData(data)
	if rainData == nil {
		t.Fatal("imgToRainData returned nil for a valid image.")
	}

	// Define the top-left corner and size of the box to print.
	// The coordinates are chosen to be near the interesting data points
	// from the TestImgToRainData test.
	xStart, yStart, boxSize := 240, 205, 20

	t.Logf("--- Rain Data Grid (section from %d,%d to %d,%d) ---", xStart, yStart, xStart+boxSize, yStart+boxSize)
	for y := yStart; y < yStart+boxSize && y < len(rainData); y++ {
		var row string
		for x := xStart; x < xStart+boxSize && x < len(rainData[y]); x++ {
			row += fmt.Sprintf("%2d ", rainData[y][x])
		}
		t.Logf("y=%-3d: %s", y, row)
	}
	t.Log("--- End of Grid ---")
}
