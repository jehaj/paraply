package main

import (
	"fmt"
	"io"
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
func (d DMIMap) UpdateMap() {
	rt := time.Now()
	t := time.Now()
	timelineRainData := make([]RainData, 6)
	for i := 0; t.Before(rt.Add(time.Hour)); i++ {
		t = t.Add(5 * time.Minute)
		rainData := getRainDataForTime(t, rt)
		timelineRainData[i] = rainData
		time.Sleep(1 * time.Second)
	}
	d.TimelineRainData = timelineRainData
}

func getRainDataForTime(t time.Time, rt time.Time) RainData {
	t = t.Truncate(5 * time.Minute)
	rt = rt.Truncate(5 * time.Minute)
	data := sendRequestForImageAt(t, rt)
	rainData := imgToRainData(data, t)
	return rainData
}

func imgToRainData(data []byte, t time.Time) RainData {
	panic("implement me")
}

func sendRequestForImageAt(t time.Time, rt time.Time) []byte {
	encodedTime := encodeTime(t)
	encodedReferenceTime := encodeTime(rt)
	url := fmt.Sprintf("https://www.dmi.dk/ZoombareKort/map?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&FORMAT=image%%2Fpng&TRANSPARENT=true&TIME=%s&REFERENCE_TIME=%s&LAYERS=nowcast_radar&WIDTH=512&HEIGHT=512&SRS=EPSG%%3A3575&BBOX=-406250%%2C-4218750%%2C250000%%2C-3562500", encodedTime, encodedReferenceTime)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:139.0) Gecko/20100101 Firefox/139.0")
	req.Header.Add("Accept", "image/avif,image/webp,image/png,image/svg+xml,image/*;q=0.8,*/*;q=0.5")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Add("Referer", "https://www.dmi.dk/")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body
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
