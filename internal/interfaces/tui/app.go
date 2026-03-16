package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/JoseGusnay/ursus/internal/application/service"
	"github.com/JoseGusnay/ursus/internal/application/usecase"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Tab represents each view in the TUI
type tab int

const (
	tabMemories tab = iota
	tabTimeline
	tabStats
)

var (
	// Colors
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	text      = lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#FAFAFA"}

	// Styles
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Background(highlight).
			Foreground(lipgloss.Color("#FFF")).
			Padding(0, 2).
			MarginRight(1)

	inactiveTabStyle = lipgloss.NewStyle().
				Background(subtle).
				Foreground(text).
				Padding(0, 2).
				MarginRight(1)

	windowStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}).
			MarginTop(1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			MarginBottom(1)
)

type item struct {
	ursus *entity.Ursus
}

func (i item) Title() string       { return i.ursus.Content }
func (i item) Description() string { return fmt.Sprintf("ID: %s | %s", i.ursus.ID[:8], i.ursus.CreatedAt.Format("2006-01-02")) }
func (i item) FilterValue() string { return i.ursus.Content }

type model struct {
	svc        *service.MemoryService
	timelineUC *usecase.GetTimelineUseCase
	list       list.Model
	input      textinput.Model
	activeTab  tab
	state      string // "VIEW", "ADD"
	err        error
	width      int
	height     int
	timeline   []usecase.TimelineDay
	stats      StatsData
}

type StatsData struct {
	TotalMemories int
	TotalSessions int
	TopTopics     []string
}

func NewModel(svc *service.MemoryService) model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Ursus Memories"
	l.SetShowStatusBar(false)
	l.SetShowTitle(false) // Use native hide method instead of style

	ti := textinput.New()
	ti.Placeholder = "Write your memory..."

	return model{
		svc:        svc,
		timelineUC: usecase.NewGetTimelineUseCase(svc.Repository()),
		list:       l,
		input:      ti,
		activeTab:  tabMemories,
		state:      "VIEW",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.fetchMemories, m.fetchTimeline, m.fetchStats)
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

func (m model) fetchTimeline() tea.Msg {
	days, err := m.timelineUC.Execute(context.Background())
	if err != nil {
		return err
	}
	return days
}

func (m model) fetchStats() tea.Msg {
	memories, _ := m.svc.List(context.Background())
	// Simple stats calculation
	topics := make(map[string]int)
	for _, memo := range memories {
		if memo.TopicKey != "" {
			topics[memo.TopicKey]++
		}
	}
	return StatsData{
		TotalMemories: len(memories),
		TopTopics:     []string{"Not implemented yet..."}, // Expand logic if needed
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "right":
			if m.state == "VIEW" {
				m.activeTab = (m.activeTab + 1) % 3
			}
		case "left":
			if m.state == "VIEW" {
				m.activeTab = (m.activeTab + 2) % 3
			}
		case "1": m.activeTab = tabMemories
		case "2": m.activeTab = tabTimeline
		case "3": m.activeTab = tabStats
		case "n":
			if m.activeTab == tabMemories && m.state == "VIEW" && !m.list.SettingFilter() {
				m.state = "ADD"
				m.input.Focus()
				return m, nil
			}
		case "esc":
			if m.state == "ADD" {
				m.state = "VIEW"
				m.input.Blur()
				return m, nil
			}
		case "enter":
			if m.state == "ADD" {
				content := m.input.Value()
				if content != "" {
					_, err := m.svc.Store(context.Background(), content, "", "", "", "")
					if err != nil {
						m.err = err
					}
					m.state = "VIEW"
					m.input.Reset()
					return m, m.fetchMemories
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(m.width-6, m.height-10)

	case []list.Item:
		m.list.SetItems(msg)
	case []usecase.TimelineDay:
		m.timeline = msg
	case StatsData:
		m.stats = msg
	case error:
		m.err = msg
	}

	if m.state == "ADD" {
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.activeTab == tabMemories {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPresione 'q' para salir.", m.err)
	}

	// Header
	tabs := []string{"[1] Memorias", "[2] Timeline", "[3] Stats"}
	header := []string{}
	for i, t := range tabs {
		if tab(i) == m.activeTab {
			header = append(header, activeTabStyle.Render(t))
		} else {
			header = append(header, inactiveTabStyle.Render(t))
		}
	}
	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, header...)

	// Main Content
	var content string
	switch m.activeTab {
	case tabMemories:
		if m.state == "ADD" {
			content = fmt.Sprintf("Nueva Memoria\n\n%s\n\n(esc para cancelar, enter para guardar)", m.input.View())
		} else {
			content = m.list.View()
		}
	case tabTimeline:
		var b strings.Builder
		b.WriteString(titleStyle.Render("Línea de Tiempo\n"))
		if len(m.timeline) == 0 {
			b.WriteString("No hay eventos registrados.")
		}
		for _, day := range m.timeline {
			b.WriteString(lipgloss.NewStyle().Foreground(highlight).Bold(true).Render("\n" + day.Date.Format("2006-01-02") + "\n"))
			for _, memo := range day.Memories {
				b.WriteString(fmt.Sprintf("  • %s: %s\n", memo.CreatedAt.Format("15:04"), memo.Content))
			}
		}
		content = b.String()
	case tabStats:
		var b strings.Builder
		contentStyle := lipgloss.NewStyle().Padding(1).Foreground(text)
		b.WriteString(titleStyle.Render("📊 Panel de Estadísticas\n"))
		
		statsGrid := lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(25).Render(fmt.Sprintf("\n📁 Total Memorias:\n   %d", m.stats.TotalMemories)),
			lipgloss.NewStyle().Width(25).Render("\n🔥 Actividad (7d):\n   Low"),
		)
		b.WriteString(contentStyle.Render(statsGrid))
		
		b.WriteString("\n\n📈 Frecuencia de Uso:\n")
		// Improved ASCII chart
		bars := []string{
			"L [▓▓░░░░░░░░] 20%",
			"M [▓▓▓▓▓▓░░░░] 60%",
			"X [▓▓▓▓▓▓▓▓░░] 80%",
			"J [▓▓▓░░░░░░░] 30%",
		}
		b.WriteString(lipgloss.NewStyle().PaddingLeft(2).Render(strings.Join(bars, "\n")))
		content = b.String()
	}

	// Footer
	footer := footerStyle.Render("Tab/flechas: Navegar • q: Salir • n: Nueva (en Memorias)")

	return lipgloss.JoinVertical(lipgloss.Left,
		headerRow,
		windowStyle.Width(m.width-4).Height(m.height-8).Render(content),
		footer,
	)
}

func Start(svc *service.MemoryService) error {
	p := tea.NewProgram(NewModel(svc), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
