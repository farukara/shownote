package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	u "github.com/farukara/shownote/utils"
)

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
                    u.Open_tasknote(taskno, c_config.notes_folder, c_config.file_ext)
            }

        case 3:
            method, taskno := args[1], args[2]
            switch method {
                case "add" , "ADD" , "a" , "A":
                    u.Add_tasknote(taskno, c_config.notes_folder, c_config.file_ext)
            }
        default:
            err := errors.New("No support for other than 2 or 3 arguments ")
            fmt.Println("usage:\nopen note for task 12: main 12\nadd note to task 12: main add 12")
            log.Error().Stack().Err(err).Msg("No support for other than 2 or 3 arguments ")
            
    }
}
