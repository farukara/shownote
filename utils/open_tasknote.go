package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func Open_tasknote (taskno, notes_folder, file_ext string) {
    task_uuid, err := Get_task_uuid(taskno)
    if err != nil {
        err := errors.New("failure to get uuid")
        log.Err(err).Stack().Str("task no", taskno).Msg( "error getting task UUID from Task Warrior" )
    }
    log.Info().Str("task uuid:", string(task_uuid))

    app := os.Getenv("EDITOR")
    homename, err := os.UserHomeDir()
    if err != nil {
        log.Err(err).Stack().Msg("failure to get user home directory")
    }
    file_name := strings.TrimSpace(string(task_uuid)) + file_ext
    arg0 := filepath.Join(homename, notes_folder,  file_name)


    // if file does not exist
    _,err = os.Stat(arg0)
    if errors.Is(err, fs.ErrNotExist) {
        log.Err(err).Stack().Msg("Note for this task does not exist.")
        fmt.Println("\n\033[7;31m==>\033[27m\033[0m Note for this task \033[1;31mdoes not\033[0m exist.")
        fmt.Println("\n1. \033[7;36mAdd one\033[0m (creates a new note for this task and annotates the task with \"Notes\")")
        fmt.Println("2. \033[7;36mCancel\033[0m")
        fmt.Print("\nChoose (1/2): ")
        s := bufio.NewScanner(os.Stdin)
        var answer string
        if s.Scan() {
            answer = s.Text()
        }
        switch answer {
            //FIX:
            case "1": 
                Add_tasknote(taskno, notes_folder, file_ext)
            case "2": 
            fmt.Println("Cancelled!")
                os.Exit(0)
            default: 
                fmt.Println("\n\033[7;31m==>\033[27m\033[0m invalid input:", answer)
                os.Exit(1)
        }
        return
    }

    cmd := exec.Command(app, arg0)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    err = cmd.Run()
    if err != nil {
        log.Err(err).Stack().Str("cmd", app + " " +arg0).Msg("error running command")
    } else {
        log.Info().Str("cmd", app + " " +arg0).Msg("note opened")
    }
}
