package editor

import (
	"gcs/view"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
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
	c := exec.Command(editor, "detail_content.env")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		fileContent, err := os.ReadFile("detail_content.env")
		equal := isEqual(secretData, fileContent)

		return EditFinishedMsg{equal, currentSecret, fileContent}
	})
}

func isEqual(secretData string, fileContent []byte) bool {
	var equal bool

	if string(fileContent) != secretData {
		log.Info().Msg("The file content is different from secretData")
		equal = false
	} else {
		log.Info().Msg("The file content is the same as secretData")
		equal = true
	}

	return equal
}
