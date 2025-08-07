package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestRequestReturnsBody(t *testing.T) {
	t1 := time.Now().UTC()
	response, _ := sendRequestForImageAt(t1, t1)
	os.WriteFile("test_temp/test_now.png", response, 0644)
	if len(response) < 10 {
		t.Log("Response length should be at least 10, but is", len(response))
		t.Fail()
	}
}

func TestRequestNextHour(t *testing.T) {
	rt := time.Now().UTC()
	time1 := time.Now().UTC()
	for i := 0; time1.Before(rt.Add(60 * time.Minute)); i++ {
		time1 = time1.Add(5 * time.Minute)
		response, _ := sendRequestForImageAt(time1, rt)
		filename := fmt.Sprintf("test_temp/test_now_%d.png", i)
		os.WriteFile(filename, response, 0644)
		time.Sleep(500 * time.Millisecond)
	}
}
