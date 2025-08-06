package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type DMIMap struct {
	TimelineRainData TimelineRainData
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
	rainData := imgToRainData(data, t)
	return rainData, nil
}

func imgToRainData(data []byte, t time.Time) RainData {
	panic("implement me")
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
	url := fmt.Sprintf("https://www.dmi.dk/ZoombareKort/map?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&FORMAT=image%%2Fpng&TRANSPARENT=true&TIME=%s&REFERENCE_TIME=%s&LAYERS=nowcast_radar&WIDTH=512&HEIGHT=512&SRS=EPSG%%3A3575&BBOX=-406250%%2C-4218750%%2C250000%%2C-3562500", encodedTime, encodedReferenceTime)

	req, err := http.NewRequest("GET", url, nil)
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

func (d DMIMap) GetPrecipitationAt(location Location) {
	//TODO implement me
	panic("implement me")
}
