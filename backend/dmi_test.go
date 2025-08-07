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

// TestGetPrecipitationAt tests the GetPrecipitationAt function using the test_future_*.png files
func TestGetPrecipitationAt(t *testing.T) {
	// Load all test_future_*.png files and convert them to RainData
	timelineRainData := make(TimelineRainData, 0)

	for i := 0; i < 12; i++ {
		testImagePath := filepath.Join("test_files", fmt.Sprintf("test_future_%d.png", i))
		data, err := os.ReadFile(testImagePath)
		if err != nil {
			t.Fatalf("Failed to read test image '%s': %v", testImagePath, err)
		}

		rainData := imgToRainData(data)
		if rainData == nil {
			t.Fatalf("imgToRainData returned nil for image %s", testImagePath)
		}

		timelineRainData = append(timelineRainData, rainData)
	}

	// Create a DMIMap with the test data
	dmiMap := &DMIMap{
		TimelineRainData: timelineRainData,
	}

	// Test cases with different locations
	testCases := []struct {
		name      string
		location  Location
		shouldErr bool
	}{
		{
			name: "Valid location in Denmark (Aarhus)",
			location: Location{
				Latitude:  56.15674,
				Longitude: 10.21076,
			},
			shouldErr: false,
		},
		{
			name: "Valid location in Denmark (Copenhagen)",
			location: Location{
				Latitude:  55.6761,
				Longitude: 12.5683,
			},
			shouldErr: false,
		},
		{
			name: "Invalid location (outside bounds)",
			location: Location{
				Latitude:  70.0,
				Longitude: 30.0,
			},
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			precipitation, err := dmiMap.GetPrecipitationAt(tc.location)

			if tc.shouldErr {
				if err == nil {
					t.Errorf("Expected error for location (%f, %f), but got none", tc.location.Latitude, tc.location.Longitude)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for location (%f, %f): %v", tc.location.Latitude, tc.location.Longitude, err)
				return
			}

			// Verify that we get precipitation data for all time points
			if len(precipitation) != len(timelineRainData) {
				t.Errorf("Expected precipitation data for %d time points, got %d", len(timelineRainData), len(precipitation))
			}

			// Verify that all precipitation values are non-negative
			for i, value := range precipitation {
				if value < 0 {
					t.Errorf("Expected non-negative precipitation value at time %d, got %d", i, value)
				}
			}

			t.Logf("Location (%f, %f) precipitation values: %v", tc.location.Latitude, tc.location.Longitude, precipitation)
		})
	}
}

// TestGetPrecipitationAt tests the GetPrecipitationAt function for Ballerup using the test_future_*.png files
func TestGetPrecipitationAtInBallerup(t *testing.T) {
	// Load all test_future_*.png files and convert them to RainData
	timelineRainData := make(TimelineRainData, 0)

	for i := 0; i < 12; i++ {
		testImagePath := filepath.Join("test_files", fmt.Sprintf("test_future_%d.png", i))
		data, err := os.ReadFile(testImagePath)
		if err != nil {
			t.Fatalf("Failed to read test image '%s': %v", testImagePath, err)
		}

		rainData := imgToRainData(data)
		if rainData == nil {
			t.Fatalf("imgToRainData returned nil for image %s", testImagePath)
		}

		timelineRainData = append(timelineRainData, rainData)
	}

	// Create a DMIMap with the test data
	dmiMap := &DMIMap{
		TimelineRainData: timelineRainData,
	}

	// Test cases with different locations
	location := Location{
		Latitude:  55.7243,
		Longitude: 12.3561,
	}

	incomingRain, err := dmiMap.GetPrecipitationAt(location)
	if err != nil {
		t.Fatalf("Unexpected error for location (%f, %f): %v", location.Latitude, location.Longitude, err)
	}
	t.Logf("Incoming rain for Ballerup: %v", incomingRain)
	expected := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 3}
	if got := incomingRain; len(got) != len(expected) {
		t.Fatalf("Expected %d time points, got %d", len(expected), len(got))
	} else {
		for i := range got {
			if got[i] != expected[i] {
				t.Errorf("At time point %d, expected %d, got %d", i, expected[i], got[i])
			}
		}
	}
}

// TestGetPrecipitationAtBoundaryConditions tests the GetPrecipitationAt function with edge cases
func TestGetPrecipitationAtBoundaryConditions(t *testing.T) {
	// Create minimal test data - just one time point
	testImagePath := filepath.Join("test_files", "test_future_0.png")
	data, err := os.ReadFile(testImagePath)
	if err != nil {
		t.Fatalf("Failed to read test image '%s': %v", testImagePath, err)
	}

	rainData := imgToRainData(data)
	if rainData == nil {
		t.Fatal("imgToRainData returned nil")
	}

	timelineRainData := TimelineRainData{rainData}
	dmiMap := &DMIMap{
		TimelineRainData: timelineRainData,
	}

	// Test with empty TimelineRainData
	emptyDmiMap := &DMIMap{
		TimelineRainData: TimelineRainData{},
	}

	validLocation := Location{
		Latitude:  56.15674,
		Longitude: 10.21076,
	}

	precipitation, err := emptyDmiMap.GetPrecipitationAt(validLocation)
	if err != nil {
		t.Errorf("Unexpected error with empty timeline data: %v", err)
	}
	if len(precipitation) != 0 {
		t.Errorf("Expected empty precipitation array, got length %d", len(precipitation))
	}

	// Test with valid location and single time point
	precipitation, err = dmiMap.GetPrecipitationAt(validLocation)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(precipitation) != 1 {
		t.Errorf("Expected precipitation array of length 1, got %d", len(precipitation))
	}
}
