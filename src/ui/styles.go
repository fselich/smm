package ui

import "github.com/charmbracelet/lipgloss"

func StyleLow() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#5a5a5a"))
}

func StyleSelected() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
		Padding(0, 0, 0, 1).
		Foreground(lipgloss.Color("#000000")).
		BorderLeftForeground(lipgloss.Color("#87CEFA")).
		Background(lipgloss.Color("#87CEFA"))
}

func StyleUnselected() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2).
		Foreground(lipgloss.Color("#87CEFA"))
}

func StyleBorder() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#87CEFA"))
}

func StyleLowBorder() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#C0C0C0"))
}

func StyleBorderTitle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Background(lipgloss.Color("#000000"))
}

func StyleToast() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Align(lipgloss.Center)
}
