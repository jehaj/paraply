package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type Map interface {
	UpdateMap()
	GetPrecipitationAt(location Location)
}

type LocationTransformer interface {
	EPSG4326To3575(latitude float64, longitude float64) (int, int)
}

type StubLocationTransformer struct{}

func (s *StubLocationTransformer) EPSG4326To3575(latitude float64, longitude float64) (int, int) {
	return 8936, -3721180
}

func makeStubLocationTransformer() *StubLocationTransformer {
	return new(StubLocationTransformer)
}

type CmdLocationTransformer struct {
	stdin  io.WriteCloser
	stdout *bufio.Reader
}

func (c *CmdLocationTransformer) EPSG4326To3575(latitude float64, longitude float64) (int, int) {
	input := fmt.Sprintf("%f %f", latitude, longitude)
	c.stdin.Write([]byte(input))
	output, _ := c.stdout.ReadString('\n')
	outputs := strings.Split(output, "\n")
	x, _ := strconv.ParseFloat(outputs[0], 64)
	y, _ := strconv.ParseFloat(outputs[1], 64)
	return int(x), int(y)
}

func makeCmdLocationTransformer() *CmdLocationTransformer {
	cmdLocationTransformer := new(CmdLocationTransformer)
	cmd := exec.Command("cs2cs", "EPSG:4326", "EPSG:3575")
	pipe, _ := cmd.StdoutPipe()
	cmdLocationTransformer.stdout = bufio.NewReader(pipe)
	cmdLocationTransformer.stdin, _ = cmd.StdinPipe()
	return cmdLocationTransformer
}

type Location struct {
	Latitude  float64
	Longitude float64
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		location := getLocationFromRequest(r)
		precipitation := getPrecipitationFor(location)
		response, _ := json.Marshal(precipitation)
		w.Write(response)
	})
	http.ListenAndServe("127.0.0.1:3000", r)
}

func getPrecipitationFor(location *Location) []int {
	return nil
}

func getLocationFromRequest(r *http.Request) *Location {
	decoder := json.NewDecoder(r.Body)
	location := new(Location)
	decoder.Decode(location)
	return location
}
