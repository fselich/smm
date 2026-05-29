package view

import (
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func CopyToClipboard(text string) tea.Cmd {
	cmds := []tea.Cmd{tea.SetClipboard(text)}

	if display := os.Getenv("DISPLAY"); display != "" {
		if _, err := exec.LookPath("xclip"); err == nil {
			cmds = append(cmds, func() tea.Msg {
				go func() {
					cmd := exec.Command("xclip", "-selection", "clipboard")
					cmd.Stdin = strings.NewReader(text)
					_ = cmd.Run()
				}()
				return nil
			})
		}
	}

	return tea.Batch(cmds...)
}
