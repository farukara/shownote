package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
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

func open_tasknote (taskno, notes_folder, file_ext string) {
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
                add_tasknote(taskno, notes_folder, file_ext)
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

func add_tasknote(taskno, notes_folder, file_ext string) {
    log.Debug().Str("taskno", taskno)
    task_uuid, err := get_task_uuid(taskno)
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
    open_tasknote(taskno, notes_folder, file_ext)
}

func init() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

}

func main() {
    type config struct {
        file_ext                        string      `yaml:"file_ext"`
        editor                          string      `yaml:"editor"`
        notes_folder                    string      `yaml:"notes_folder"`
        open_tasknote_after_creation    int         `yaml:"open_tasknote_after_creation"`
    }

    c_config := config {
        file_ext: ".md",
        editor: "vim",
        notes_folder: ".tasknotes",
        open_tasknote_after_creation: 1,
    }// current config

    // NOTE: add home folder option for config file
    // NOTE: add backup for no file cases
    // NOTE: move to its own function
    homename, err := os.UserHomeDir()
    if err != nil {
        log.Err(err).Stack().Msg("failure to get user home directory")
    }
    configpath := ""
    configpath1 := filepath.Join(homename, ".config/shownote/config.yaml")
    configpath2 := filepath.Join(homename, ".shownote/config.yaml")
    _,err1 := os.Stat(configpath1)
    _,err2 := os.Stat(configpath2)
    if errors.Is(err2, fs.ErrNotExist) {
        log.Info().Msg("config.yaml is not at $HOME/.shownote folder")
    } else {
        configpath = configpath2
    }
    if errors.Is(err1, fs.ErrNotExist) {
        log.Info().Msg("config.yaml is not at $HOME/.config/.shownote folder")
    } else {
        configpath = configpath1
    }
    
    log.Debug().Str("len of configpath", strconv.Itoa(len(configpath)))
    if len(configpath) != 0 {
        config_file,err := os.Open(configpath)
        defer config_file.Close()
        b_config, err := io.ReadAll(config_file)
        if err != nil {
            log.Err(err).Stack().Msg("error reading bytes from config file")
        }
        err = yaml.Unmarshal(b_config, &c_config)
        if err != nil {
            log.Err(err).Stack().Msg("error unmarshalling config file")
        }
    }

    args := os.Args

    switch len(os.Args) {
        case 2:
            switch args[1] {
                case "tidy", "t":
                    // NOTE: implement
                default:
                    taskno := args[1]
                    open_tasknote(taskno, c_config.notes_folder, c_config.file_ext)
            }

        case 3:
            method, taskno := args[1], args[2]
            switch method {
                case "add" , "ADD" , "a" , "A":
                    add_tasknote(taskno, c_config.notes_folder, c_config.file_ext)
            }
        default:
            err := errors.New("No support for other than 2 or 3 arguments ")
            fmt.Println("usage:\nopen note for task 12: main 12\nadd note to task 12: main add 12")
            log.Error().Stack().Err(err).Msg("No support for other than 2 or 3 arguments ")
            
    }
}
