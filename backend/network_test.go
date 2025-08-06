package main

import (
	"os"
	"testing"
	"time"
)

func TestRequestReturnsBody(t *testing.T) {
	t1 := time.Now()
	response := sendRequestForImageAt(t1, t1)
	os.WriteFile("test_gen/test_now.png", response, 0644)
	if len(response) < 10 {
		t.Log("Response length should be at least 10, but is", len(response))
		t.Fail()
	}
}
