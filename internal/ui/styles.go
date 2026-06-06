package ui

import (
	lipgloss "charm.land/lipgloss/v2"
)

var lightDark = lipgloss.LightDark(true)

func StyleLow() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#5a5a5a"))
}

func StyleSelected() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lightDark(lipgloss.Color("#F793FF"), lipgloss.Color("#AD58B4"))).
		Foreground(lipgloss.Color("#EE6FF8")).
		Padding(0, 0, 0, 1).
		Foreground(lipgloss.Color("#000000")).
		BorderLeftForeground(lipgloss.Color("#87CEFA")).
		Background(lipgloss.Color("#87CEFA"))
}

func StyleUnselected() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lightDark(lipgloss.Color("#1a1a1a"), lipgloss.Color("#dddddd"))).
		Padding(0, 0, 0, 2).
		Foreground(lipgloss.Color("#87CEFA"))
}

func StyleBorder(selected bool) lipgloss.Style {
	var color string
	if selected {
		color = "#87CEFA"
	} else {
		color = "#505050"
	}
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(color))
}

func StyleModal() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Border(lipgloss.DoubleBorder()).
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

func StyleSelectedButton() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#87CEFA")).
		Padding(0, 1)
}

func StyleUnselectedButton() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)
}

func StyleDisabledButton() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555")).
		Padding(0, 1)
}
