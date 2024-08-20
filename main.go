package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/dionvu/goplay/player"
	kb "github.com/eiannone/keyboard"
	"github.com/gopxl/beep/speaker"
)

func main() {
	kb.Open()

	defer kb.Close()

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	player, err := player.New(f)
	if err != nil {
		log.Fatal(err)
	}
	defer player.Streamer.Close()

	speaker.Play(player.Resampler)

	for {
		char, key, _ := kb.GetKey()

		switch char {
		case 'p':
			player.VolumeUp(0.1)
			fmt.Println(math.Round(player.Volume()*10) / 10)
		}

		switch key {
		case kb.KeyCtrlC:
			return
		}
	}
}
