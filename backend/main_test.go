package main

import (
	"testing"
	"time"
)

type StubLocationTransformer struct{}

func (s *StubLocationTransformer) EPSG4326To3575(latitude float64, longitude float64) (int, int) {
	return 8936, -3721180
}

func makeStubLocationTransformer() *StubLocationTransformer {
	return new(StubLocationTransformer)
}

func TestExample(t *testing.T) {
	t.Log("It is very much working.")
}

func TestProj(t *testing.T) {
	projTransformer := makeProjTransformer()
	lat := 56.15674
	lon := 10.21076
	x, y := projTransformer.EPSG4326To3575(lat, lon)
	if x != 13689 && y != -3721342 {
		t.Log("Expected 13689, -3721342, got", x, y)
		t.Fail()
	}
}

func TestCmd(t *testing.T) {
	t.Log("Creating another process...")
	cmdTransformer := makeCmdLocationTransformer()
	lat := 56.15674
	lon := 10.21076
	t.Log("Transforming coordinate")
	x, y := cmdTransformer.EPSG4326To3575(lat, lon)
	if x != 13689 && y != -3721342 {
		t.Log("Expected 13689, -3721342, got", x, y)
		t.Fail()
	}
}

func TestEncodingOfTime(t *testing.T) {
	t1 := time.Date(2025, time.July, 7, 12, 20, 0, 0, time.UTC)
	encodedTime := encodeTime(t1)
	if encodedTime != "2025-07-07T12%3A20%3A00Z" {
		t.Logf("Expected 2025-07-07T12%%3A20%%3A00Z, got %s", encodedTime)
		t.Fail()
	}
}
