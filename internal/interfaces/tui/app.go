package tui

import (
	"context"
	"fmt"

	"github.com/JoseGusnay/ursus/internal/application/service"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)
	itemStyle = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#AD58B4"))
)

type item struct {
	ursus *entity.Ursus
}

func (i item) Title() string       { return i.ursus.Content }
func (i item) Description() string { return fmt.Sprintf("ID: %s | %s", i.ursus.ID[:8], i.ursus.CreatedAt.Format("2006-01-02")) }
func (i item) FilterValue() string { return i.ursus.Content }

type model struct {
	list     list.Model
	svc      *service.MemoryService
	err      error
	state    string // "LIST", "ADD", "SEARCH"
	input    textinput.Model
	width    int
	height   int
}

func NewModel(svc *service.MemoryService) model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Ursus Memories"
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle

	ti := textinput.New()
	ti.Placeholder = "Write your memory..."
	ti.Focus()

	return model{
		list:  l,
		svc:   svc,
		state: "LIST",
		input: ti,
	}
}

func (m model) Init() tea.Cmd {
	return m.fetchMemories
}

func (m model) fetchMemories() tea.Msg {
	results, err := m.svc.List(context.Background())
	if err != nil {
		return err
	}
	items := make([]list.Item, len(results))
	for i, r := range results {
		items[i] = item{ursus: r}
	}
	return items
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == "LIST" && !m.list.SettingFilter() {
				return m, tea.Quit
			}
		case "n":
			if m.state == "LIST" && !m.list.SettingFilter() {
				m.state = "ADD"
				m.input.Reset()
				m.input.Focus()
				return m, nil
			}
		case "esc":
			if m.state == "ADD" {
				m.state = "LIST"
				return m, nil
			}
		case "enter":
			if m.state == "ADD" {
				content := m.input.Value()
				if content != "" {
					_, err := m.svc.Store(context.Background(), content, "")
					if err != nil {
						m.err = err
					}
					m.state = "LIST"
					return m, m.fetchMemories
				}
			}
		case "x": // Delete
			if m.state == "LIST" && !m.list.SettingFilter() {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					err := m.svc.Delete(context.Background(), i.ursus.ID)
					if err != nil {
						m.err = err
					}
					return m, m.fetchMemories
				}
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.width = msg.Width
		m.height = msg.Height

	case []list.Item:
		m.list.SetItems(msg)
		return m, nil

	case error:
		m.err = msg
		return m, nil
	}

	if m.state == "ADD" {
		m.input, cmd = m.input.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPresiona 'q' para salir.", m.err)
	}

	switch m.state {
	case "ADD":
		return docStyle.Render(
			fmt.Sprintf(
				"Add New Memory\n\n%s\n\n(esc to cancel, enter to save)",
				m.input.View(),
			),
		)
	default:
		return docStyle.Render(m.list.View())
	}
}

func Start(svc *service.MemoryService) error {
	p := tea.NewProgram(NewModel(svc), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
