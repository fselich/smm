package editor

import (
	"os"
	"os/exec"
	"path/filepath"
	"smm/internal/view"

	tea "github.com/charmbracelet/bubbletea"
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

	tempDir := os.TempDir()
	hash := currentSecret.Hash()
	filePath := filepath.Join(tempDir, hash)

	c := exec.Command(editor, filePath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		fileContent, err := os.ReadFile(filePath)
		equal := string(fileContent) == secretData

		_ = os.Remove(filePath)

		return EditFinishedMsg{equal, currentSecret, fileContent}
	})
}
