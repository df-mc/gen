package main

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/gen"
	"github.com/ojrac/opensimplex-go"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
)

import _ "net/http/pprof"

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	srv := server.New(&config, log)
	srv.CloseOnProgramEnd()
	if err := srv.Start(); err != nil {
		log.Fatalln(err)
	}
	srv.World().Generator(gen.New())
	srv.World().ReadOnly()
	srv.World().SetTime(5000)
	srv.World().StopTime()

	/*const (
		dx, dy = 10000, 10000
		warp = 70
		dd = 0.01
	)
	i := image.NewRGBA(image.Rect(0, 0, 10000, 10000))
	n := noise(3, 4, 2, 0.5)
	max := func(x float64) float64 {
		if x < 0 {
			return 0
		} else if x > 1 {
			return 1
		}
		return x
	}
	for xx := 0; xx < dx; xx++ {
		for yy := 0; yy < dy; yy++ {
			x, y := float64(xx), float64(yy)
			warpX, warpY := n(x*dd, y*dd), n(x*dd+5.2, y*dd+1.3)
			warpX, warpY = n(x*dd+warpX*warp+1.7, y*dd+warpY*warp+9.2), n(x*dd+warpX*warp+8.3, y*dd+warpY*warp+2.8)
			i.Set(xx, yy, color.RGBA{
				A: uint8(max(n(x*dd+warpX*warp, y*dd+warpY*warp)+0.5) * 255),
			})
		}
		if xx % 1000 == 0 {
			fmt.Println(xx)
		}
	}
	f, _ := os.OpenFile("image.png", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err := png.Encode(f, i); err != nil {
		panic(err)
	}

	return*/
	for {
		if _, err := srv.Accept(); err != nil {
			return
		}
	}
}

// readConfig reads the configuration from the config.toml file, or creates the file if it does not yet exist.
func readConfig() (server.Config, error) {
	c := server.DefaultConfig()
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := ioutil.WriteFile("config.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}

func noise(seed int64, octaves int, lacunarity, persistence float64) func(x, y float64) float64 {
	// Create noise instances for all octaves.
	n := make([]opensimplex.Noise, octaves)
	for i := 0; i < octaves; i++ {
		n[i] = opensimplex.New(seed)
	}

	// Calculate the maximum possible value when adding together noise values from all octaves.
	var max float64
	for i := 0; i < octaves; i++ {
		max += math.Pow(persistence, float64(i))
	}

	return func(x, y float64) float64 {
		var v float64

		for i, noise := range n {
			// Lacunarity gives us the frequency
			freq := math.Pow(lacunarity, float64(i)) * 0.04
			amp := math.Pow(persistence, float64(i))

			v += amp * noise.Eval2(x*freq, y*freq)
		}
		// Normalise at the end so we get a value in the range [0-1].
		v /= max
		return v
	}
}
