package editor

import (
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
	"smm/view"
)

type EditFinishedMsg struct {
	Equal         bool
	CurrentSecret view.Secret
	SecretData    []byte
}

func OpenEditor(secretData string, currentSecret view.Secret) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	hash := currentSecret.Hash()
	c := exec.Command(editor, "/tmp/"+hash)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		fileContent, err := os.ReadFile("/tmp/" + hash)
		equal := string(fileContent) == secretData

		_ = os.Remove("/tmp/" + hash)

		return EditFinishedMsg{equal, currentSecret, fileContent}
	})
}
