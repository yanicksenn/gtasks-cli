package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	unfocusedSelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("240"))
)

type itemDelegate struct{
	focused bool
}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(taskListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Title())

	fn := itemStyle.Render
	if index == m.Index() {
		if d.focused {
			fn = func(s ...string) string {
				return selectedItemStyle.Render("> " + s[0])
			}
		} else {
			fn = func(s ...string) string {
				return unfocusedSelectedItemStyle.Render("  " + s[0])
			}
		}
	}

	fmt.Fprint(w, fn(str))
}

type taskItemDelegate struct{}

func (d taskItemDelegate) Height() int                               { return 1 }
func (d taskItemDelegate) Spacing() int                              { return 0 }
func (d taskItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d taskItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(taskItem)
	if !ok {
		return
	}

	status := "[ ]"
	if i.Status == "completed" {
		status = "[x]"
	}

	str := fmt.Sprintf("%s %s", status, i.Title())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}
