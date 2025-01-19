package dialog

import (
	"github.com/charmbracelet/lipgloss"
	"pirate-wars/cmd/screen"
)

var BorderStyle = lipgloss.Border{
	Top:         "─",
	Bottom:      "─",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "└",
	BottomRight: "┘",
}

func SetScreenStyle(width int, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("0")).
		Width(width).
		Height(height)
}

func GetSidebarStyle() lipgloss.Style {
	var SidebarWidth = (screen.InfoPaneSize * 3)
	if SidebarWidth > 25 {
		SidebarWidth += 1
	} else if SidebarWidth > 18 {
		SidebarWidth += 2
	} else {
		SidebarWidth += 3
	}
	return lipgloss.NewStyle().
		Align(lipgloss.Left).
		Border(BorderStyle).
		Foreground(lipgloss.Color("#FAFAFA")).
		BorderForeground(lipgloss.Color("33")).
		BorderBackground(lipgloss.Color("0")).
		Background(lipgloss.Color("0")).
		//Margin(1, 1, 0, 0).
		Padding(1).
		Height(20).
		Width(SidebarWidth)
}

var base = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("0"))

var bullet = lipgloss.NewStyle().SetString("·").
	Foreground(lipgloss.Color("#43BF6D")).
	Background(lipgloss.Color("0")).
	PaddingRight(1).
	String()

var ListItem = func(s string) string {
	return bullet + lipgloss.NewStyle().
		Foreground(lipgloss.Color("#969B86")).
		Background(lipgloss.Color("0")).
		Render(s)
}

var ListHeader = base.
	//BorderStyle(lipgloss.NormalBorder()).
	//BorderBottom(true).
	Background(lipgloss.Color("0")).
	BorderBackground(lipgloss.Color("0")).
	PaddingBottom(1).
	Width(100).
	Render
