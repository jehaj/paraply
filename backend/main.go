package main

import (
	"bufio"
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
	GetRainNiveauAt(location Location)
}

type LocationTransformer interface {
	EPSG4326To3575(latitude float64, longitude float64) (int, int)
}

type CmdLocationTransformer struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

func (c *CmdLocationTransformer) EPSG4326To3575(latitude float64, longitude float64) (int, int) {
	input := fmt.Sprintf("%f %f", latitude, longitude)
	c.stdin.Write([]byte(input))
	bufferedReader := bufio.NewReader(c.stdout)
	output, _ := bufferedReader.ReadString('\n')
	outputs := strings.Split(output, "\n")
	x, _ := strconv.Atoi(outputs[0])
	y, _ := strconv.Atoi(outputs[1])
	return x, y
}

func makeCmdLocationTransformer() *CmdLocationTransformer {
	cmdLocationTransformer := new(CmdLocationTransformer)
	cmd := exec.Command("cs2cs", "EPSG:4326", "EPSG:3575")
	cmdLocationTransformer.stdout, _ = cmd.StdoutPipe()
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
	http.ListenAndServe("127.0.0.1:3000", r)
}
