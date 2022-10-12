package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func get_task_uuid(taskno string) ([]byte, error){
    app := "task"
    arg0 := "_get"
    arg1 := taskno + ".uuid"

    cmd := exec.Command(app, arg0, arg1)
    // NOTE: find out what is returned by err from cmd
    task_uuid, err := cmd.Output()
    if err == nil && len(task_uuid) == 0 {
        err = errors.New("no task with that number")
    }
    return task_uuid, err
}

func open_tasknote (taskno string) {
    task_uuid, err := get_task_uuid(taskno)
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
    arg0 := homename + "/.tasknotes/" + strings.TrimSpace(string(task_uuid)) + ".md"

    //check if file exists
    _,err = os.Stat(arg0)
    if errors.Is(err, fs.ErrNotExist) {
        log.Err(err).Stack().Msg("Note for this task does not exist.")
        fmt.Println("\n1. Add one (creates a new note for this task and annotates the task with \"Notes\")")
        fmt.Println("2. Cancel")
        fmt.Println("Choose (1/2): ")
        s := bufio.NewScanner(os.Stdin)
        var answer string
        if s.Scan() {
            answer = s.Text()
        }
        switch answer {
            //FIX:
            case "1": 
                
            case "2": 
                os.Exit(0)
            default: 
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
    }
}

func add_tasknote(taskno, method string) {
    log.Debug().Str("taskno", taskno)
    log.Debug().Str("method", method)
    task_uuid, err := get_task_uuid(taskno)
    if err != nil {
        err := errors.New("failure to get uuid")
        log.Err(err).Stack().Str("task no", taskno).Msg( "error getting task UUID from Task Warrior" )
    }
    if method != "add" {
        err := errors.New("unsopported method")
        log.Err(err).Stack().Str("unsported method", method)
    }
    homename, err := os.UserHomeDir()
    if err != nil {
        log.Err(err).Stack().Msg("failure to get user home directory")
    }
    log.Info().Str("task uuid:", string(task_uuid))
    // filename := homename + "/.tasknotes/" + strings.TrimSpace(string(task_uuid)) + ".md"
    filename := homename + "/Documents/golang/tasknotes/" + strings.TrimSpace(string(task_uuid)) + ".md"
    filePath, _ := filepath.Abs(filename)
    os.Create(filePath)
    log.Info().Str("file_name", filename).Msg("file created")
    open_tasknote(taskno)
}

func init() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
    args := os.Args

    switch len(os.Args) {
        case 2:
            taskno := args[1]
            open_tasknote(taskno)

        case 3:
            method, taskno := args[1], args[2]
            add_tasknote(taskno, method)
        default:
            err := errors.New("No support for other than 2 or 3 arguments ")
            fmt.Println("usage:\nopen note for task 12: main 12\nadd note to task 12: main add 12")
            log.Error().Stack().Err(err).Msg("No support for other than 2 or 3 arguments ")
            
    }
}

