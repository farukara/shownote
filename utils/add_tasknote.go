package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func Add_tasknote(taskno, notes_folder, file_ext string) {
    log.Debug().Str("taskno", taskno)
    task_uuid, err := Get_task_uuid(taskno)
    if err != nil {
        err := errors.New("failure to get uuid")
        log.Err(err).Stack().Str("task no", taskno).Msg( "error getting task UUID from Task Warrior" )
    }

    homename, err := os.UserHomeDir()
    if err != nil {
        log.Err(err).Stack().Msg("failure to get user home directory")
    }
    log.Info().Str("task uuid:", string(task_uuid))
    note_path := filepath.Join(homename, notes_folder,  strings.TrimSpace(string(task_uuid)) + file_ext)
    os.Create(note_path)
    log.Info().Str("file_name", filepath.Base(note_path) ).Msg("file created")
    fmt.Println("\n\033[7;32m==>\033[27m\033[0m a note added to task", taskno)
    Open_tasknote(taskno, notes_folder, file_ext)
}
