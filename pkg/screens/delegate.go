package screens

import (
	"fmt"
	"io"

	"github.com/K44Z/kzchat/internal/server/schemas"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)

	unselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
)

type UserDelegate struct{}

func (d UserDelegate) Height() int                               { return 1 }
func (d UserDelegate) Spacing() int                              { return 1 }
func (d UserDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d UserDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	user, ok := listItem.(schemas.User)
	if !ok {
		return
	}
	if index == m.Index() {
		fmt.Fprintf(w, "%s%s",
			m.Styles.FilterCursor.Render("> "),
			selectedStyle.Render(user.Username))
	} else {
		fmt.Fprintf(w, "  %s", unselectedStyle.Render(user.Username))
	}
}
