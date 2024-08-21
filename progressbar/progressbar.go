package progressbar

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	pb "github.com/cheggaaa/pb/v3"
	"github.com/dionvu/gomp/player"
)

const (
	forwardInterval  = 5
	BackwardInterval = 5
)

type tickMsg struct{}

type Model struct {
	bar    *pb.ProgressBar
	player *player.Player
}

func New(total int, player *player.Player) Model {
	bar := pb.StartNew(int(player.Total().Seconds()))
	bar.SetMaxWidth(80)

	return Model{
		player: player,
		bar:    bar,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.bar.Current() < m.bar.Total() {
			m.bar.Current()
			m.bar.Increment()
			return m, tea.Tick(time.Second, func(time.Time) tea.Msg {
				return tickMsg{}
			})
		}

		m.bar.SetCurrent(m.bar.Total())
		m.bar.Finish()

		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "f":
			m.bar.Add(forwardInterval)

			if m.bar.Current() >= m.bar.Total() {
				m.bar.SetCurrent(m.bar.Total())

				m.bar.Finish()
			}

			m.player.Forward(forwardInterval)

		case "b":
			m.bar.Add(-BackwardInterval)

			if m.bar.Current() <= 0 {
				m.bar.SetCurrent(0)
			}

			m.player.Backward(BackwardInterval)
		}
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) View() string {
	return ""
}
