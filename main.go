package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/gomp/player"
	pb "github.com/dionvu/gomp/progressbar"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no file in argument")
		return
	}

	arg := os.Args[1]

	f, err := os.Open(arg)
	if err != nil {
		fmt.Println("invalid file, `gomp help` for help")
		return
	}

	player, err := player.New(f)
	if err != nil {
		log.Fatal(err)
	}
	defer player.Close()

	player.Start()

	tp := tea.NewProgram(pb.New(10, player))
	if _, err := tp.Run(); err != nil {
		log.Fatal(err)
	}
}
