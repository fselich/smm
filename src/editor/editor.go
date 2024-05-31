package editor

import (
	"gcs/view"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
)

type EditorFinishedMsg struct{ Err error }

func OpenEditor(secretData string, currentSecret view.Secret) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, "detail_content.env")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		compareAndDeleteFile(secretData, "detail_content.env")

		return EditorFinishedMsg{err}
	})
}

func compareAndDeleteFile(secretData string, file string) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		log.Fatal()
	}

	if string(fileContent) != secretData {
		log.Info().Msg("The file content is different from secretData")
	} else {
		log.Info().Msg("The file content is the same as secretData")
	}

	err = os.Remove(file)
	if err != nil {
		log.Fatal()
	}
}
