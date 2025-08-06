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

func (d DMIMap) UpdateMap() {
	//TODO implement me
	panic("implement me")
}

func getRainDataForTime(t time.Time) RainData {
	panic("what to use as reference time")
	data := sendRequestForImageAt(t, t)
	rainData := imgToRainData(data)
	return rainData
}

func imgToRainData(data []byte) RainData {
	panic("implement me")
}

func sendRequestForImageAt(t time.Time, rt time.Time) []byte {
	t = t.Truncate(5 * time.Minute)
	rt = rt.Truncate(5 * time.Minute)
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
