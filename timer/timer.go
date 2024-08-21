package timer

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/gomp/player"
)

const (
	forwardInterval  = 5
	BackwardInterval = 5
)

type tickMsg struct{}

type Model struct {
	player *player.Player
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func New(total int, player *player.Player) Model {
	return Model{
		player: player,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.player.Current().Seconds() < m.player.Total().Seconds() {
			return m, tea.Tick(time.Second, func(time.Time) tea.Msg {
				return tickMsg{}
			})
		}

		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "f":
			m.player.Forward(forwardInterval)

		case "b":

			m.player.Backward(BackwardInterval)

		case "u":
			m.player.VolumeUp(0.1)

		case "d":
			m.player.VolumeDown(0.1)

		case "p", " ":
			m.player.PlayPause()
		}
	}

	return m, nil
}

func (m Model) View() string {
	curr := secondsFormat(int(m.player.Current().Seconds()))

	header := fmt.Sprintf("%10s %-10s %-10s Volume: %v", "", curr,
		secondsFormat(int(m.player.Total().Seconds())), m.player.Volume())

	help := fmt.Sprintf("%10s %s %-10s", "", "p", "pause") +
		fmt.Sprintf("%s %-15s", "f", fmt.Sprint("foward ", forwardInterval, "s")) +
		fmt.Sprintf("%s %s\n", "b", fmt.Sprint("backward ", BackwardInterval, "s"))

	help += fmt.Sprintf("%10s %s %-10s", "", "q", "quit") +
		fmt.Sprintf("%s %-15s", "b", fmt.Sprint("back ", BackwardInterval, "s")) +
		fmt.Sprintf("%s %s\n", "b", fmt.Sprint("backward ", BackwardInterval, "s"))

	return "\n" + header + "\n\n" + help + "\n"
}

func secondsFormat(seconds int) string {
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	return fmt.Sprintf("[%vh %vm %vs]", hours, minutes, seconds)
}
