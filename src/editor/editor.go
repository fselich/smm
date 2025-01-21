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
	hash := currentSecret.Hash()
	c := exec.Command(editor, "/tmp/"+hash)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		fileContent, err := os.ReadFile("/tmp/" + hash)
		equal := isEqual(secretData, fileContent)

		_ = os.Remove("/tmp/" + hash)

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
