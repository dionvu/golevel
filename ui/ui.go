package ui

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

		case "f", "l":
			if m.player.Current().Seconds() >= m.player.Total().Seconds() {
				return m, tea.Quit
			}

			m.player.Forward(forwardInterval)

		case "b", "h":
			m.player.Backward(BackwardInterval)

		case "u", "k":
			m.player.VolumeUp(0.1)

		case "d", "j":
			m.player.VolumeDown(0.1)

		case "p", " ":
			m.player.PlayPause()
		}
	}

	return m, nil
}

func (m Model) View() string {
	curr := secondsFormat(int(m.player.Current().Seconds()))
	total := secondsFormat(int(m.player.Total().Seconds()))
	volume := m.player.Volume()

	header := fmt.Sprintf(
		"%5s %-10s %-10s Volume: %.1f",
		"",
		curr,
		total,
		volume,
	)

	help := fmt.Sprintf(
		"%5s %s %-15s %s %-15s %s %s\n",
		"", // Ident
		"f", fmt.Sprintf("forward %ds", forwardInterval),
		"u", "volume up",
		"p", "pause",
	) + fmt.Sprintf(
		"%5s %s %-15s %s %-15s %s %s\n",
		"", // Ident
		"b", fmt.Sprintf("back %ds", BackwardInterval),
		"d", "volume down",
		"q", "quit",
	)

	return "\n" + header + "\n\n" + help
}

func secondsFormat(seconds int) string {
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	return fmt.Sprintf("[%vh %vm %vs]", hours, minutes, seconds)
}
