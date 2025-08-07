package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type DMIMap struct {
	TimelineRainData TimelineRainData
}

type DMIBounds struct {
	left, bottom, right, top int
}

func getBounds() DMIBounds {
	return DMIBounds{
		left:   -406250,
		bottom: -4218750,
		right:  250000,
		top:    -3562500,
	}
}

// UpdateMap gets new data from DMI and updates the TimelineRainData.
// To minimise the risk of getting 429 Too Many Requests, we wait one
// second between every request. This also means that this function
// is blocking. Do not call on main thread!
// TODO: use mutex when updating d.TimelineRainData! since it might
// be done from another thread.
func (d *DMIMap) UpdateMap() {
	rt := time.Now()
	// A forecast for the next hour in 5-minute steps means 12 data points.
	timelineRainData := make(TimelineRainData, 12)
	for i := 0; i < 12; i++ {
		// Calculate forecast time, starting from 5 minutes from now.
		forecastTime := rt.Add(time.Duration(i+1) * 5 * time.Minute)
		rainData, err := getRainDataForTime(forecastTime, rt)
		if err != nil {
			// Log the error and continue. In a real application, you might want
			// more sophisticated error handling, like retries or circuit breakers.
			fmt.Printf("failed to get rain data for %v: %v\n", forecastTime, err)
			continue
		}
		timelineRainData[i] = rainData
		time.Sleep(1 * time.Second)
	}
	d.TimelineRainData = timelineRainData
}

// getRainDataForTime fetches a weather radar image from DMI for a specific time,
// processes it, and returns it as RainData.
// The time `t` is the forecast time, and `rt` is the reference time (the time
// the forecast was made).
func getRainDataForTime(t time.Time, rt time.Time) (RainData, error) {
	data, err := sendRequestForImageAt(t, rt)
	if err != nil {
		return nil, fmt.Errorf("could not get image data for time %v: %w", t, err)
	}
	rainData := imgToRainData(data)
	return rainData, nil
}

type Color struct {
	R, G, B int
}

func imgToRainData(data []byte) RainData {
	rainColors := []Color{
		{92, 0, 51},
		{128, 0, 0},
		{204, 31, 31},
		{230, 57, 57},
		{255, 82, 82},
		{255, 124, 124},
		{255, 181, 181},
		{255, 142, 82},
		{255, 178, 0},
		{255, 217, 0},
		{0, 143, 233},
		{61, 171, 238},
		{109, 191, 242},
		{22, 225, 204},
		{125, 238, 226},
		{158, 242, 233},
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("failed to decode image: %v", err)
		return nil
	}

	// Create a 2D slice to hold the rain intensity data, matching the image dimensions.
	imgHeight := img.Bounds().Dy()
	imgWidth := img.Bounds().Dx()

	const colorMatchMargin = 40

	rainData := make(RainData, imgHeight)

	// Iterate over each pixel of the image.
	for y := 0; y < imgHeight; y++ {
		rainData[y] = make([]int, imgWidth)
		for x := 0; x < imgWidth; x++ {
			// Get the color of the pixel. Note that RGBA() returns values in the
			// range [0, 65535], not [0, 255], so we need to scale them.
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// Compare the pixel color with the predefined rain colors, iterating backwards
			// to check for lighter rain first.
			// TODO make different choosing algorithms... this works badly. Check which
			// color it is closest to and choose that.
			numberOfColors := len(rainColors)
			for i, color := range rainColors {
				// If a match is found within a certain margin, store the color's index.
				// This accounts for minor color variations from compression artifacts.
				if abs(int(r8)-color.R) <= colorMatchMargin &&
					abs(int(g8)-color.G) <= colorMatchMargin &&
					abs(int(b8)-color.B) <= colorMatchMargin {
					rainData[y][x] = numberOfColors - i
					break
				}
			}
		}
	}

	return rainData
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// sendRequestForImageAt constructs and sends a GET request to the DMI weather
// radar service to fetch a map image for a specific time. It mimics a browser
// request to avoid being blocked.
//
// - IMPORTANT! The time.Time objects must have location time.UTC! (TODO)
//
// The `t` parameter is the forecast time for the image, and `rt` is the
// reference time (when the forecast was made). Both are truncated to 5 minutes.
// before being used in the request.
// It returns the raw byte data of the PNG image, or an error if the request fails.
func sendRequestForImageAt(t time.Time, rt time.Time) ([]byte, error) {
	t = t.Truncate(5 * time.Minute)
	rt = rt.Truncate(5 * time.Minute)
	log.Printf("Requesting image for time %s, reference time %s\n", t.Format(time.RFC3339), rt.Format(time.RFC3339))
	encodedTime := encodeTime(t)
	encodedReferenceTime := encodeTime(rt)
	bounds := getBounds()
	urlFormatted := fmt.Sprintf("https://www.dmi.dk/ZoombareKort/map?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&FORMAT=image%%2Fpng&TRANSPARENT=true&"+
		"TIME=%s&"+
		"REFERENCE_TIME=%s&LAYERS=nowcast_radar&WIDTH=512&HEIGHT=512&SRS=EPSG%%3A3575&"+
		"BBOX=%d%%2C%d%%2C%d%%2C%d",
		encodedTime,
		encodedReferenceTime,
		bounds.left,
		bounds.bottom,
		bounds.right,
		bounds.top)

	req, err := http.NewRequest("GET", urlFormatted, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:139.0) Gecko/20100101 Firefox/139.0")
	req.Header.Add("Accept", "image/avif,image/webp,image/png,image/svg+xml,image/*;q=0.8,*/*;q=0.5")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Add("Referer", "https://www.dmi.dk/")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// encodeTime takes a time.Time t and returns it as a URL-encoded RFC3339 string.
func encodeTime(t time.Time) string {
	formattedTime := t.Format(time.RFC3339)
	encodedTime := url.QueryEscape(formattedTime)
	return encodedTime
}

func (d *DMIMap) GetPrecipitationAt(location Location) {
	//TODO implement me
	panic("implement me")
}
