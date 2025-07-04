package views

// A simple example that shows how to send activity to Bubble Tea in real-time
// through a channel.

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/optique-dev/optique"
)

// A message used to indicate that activity has occurred. In the real world (for
// example, chat) this would contain actual data.
type responseMsg struct{}

// Simulate a process that sends events at an irregular interval in real time.
// In this case, we'll send events on the channel at a random interval between
// 100 to 1000 milliseconds. As a command, Bubble Tea will run this
// asynchronously.
func listenForActivity(sub chan struct{}, cmd *exec.Cmd) tea.Cmd {
	return func() tea.Msg {
		if err := cmd.Run(); err != nil {
			optique.Error(fmt.Sprintf("error running command: %s", err))
			os.Exit(1)
		}
		sub <- struct{}{}
		return responseMsg{}
	}
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

type model struct {
	label     string
	cmd       *exec.Cmd
	sub       chan struct{} // where we'll receive activity notifications
	responses int           // how many responses we've received
	spinner   spinner.Model
	quitting  bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		listenForActivity(m.sub, m.cmd), // generate activity
		waitForActivity(m.sub),          // wait for activity
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case responseMsg:
		m.quitting = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	s := fmt.Sprintf("\n %s %s\n", m.spinner.View(), m.label)
	return s
}

func Load(cmd *exec.Cmd, label string) {
	p := tea.NewProgram(model{
		sub:     make(chan struct{}),
		cmd:     cmd,
		label:   label,
		spinner: spinner.New(),
	})

	if _, err := p.Run(); err != nil {
		optique.Error(fmt.Sprintf("could not start program: %s", err))
		os.Exit(1)
	}
}
