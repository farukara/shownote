package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	u "github.com/farukara/shownote/utils"
)

type config struct {
    file_ext                        string      `yaml:"file_ext"`
    editor                          string      `yaml:"editor"`
    notes_folder                    string      `yaml:"notes_folder"`
    open_tasknote_after_creation    int         `yaml:"open_tasknote_after_creation"`
}

func init() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
    // current config
    c_config := config {
        file_ext: ".md",
        editor: "vim",
        notes_folder: ".tasknotes",
        open_tasknote_after_creation: 1,
    }

    // NOTE: add home folder option for config file
    // NOTE: add backup for no file cases
    // NOTE: move to its own function
    homename, err := os.UserHomeDir()
    if err != nil {
        log.Err(err).Stack().Msg("failure to get user home directory")
    }
    is_config_exist := false
    configpath := ""
    configpath1 := filepath.Join(homename, ".config/shownote/config.yaml")
    configpath2 := filepath.Join(homename, ".shownote/config.yaml")
    _,err1 := os.Stat(configpath1)
    _,err2 := os.Stat(configpath2)
    if !errors.Is(err2, fs.ErrNotExist) {
        configpath = configpath2
        is_config_exist = true
    }
    if !errors.Is(err1, fs.ErrNotExist) {
        configpath = configpath1
        is_config_exist = true
    }
    if !is_config_exist {
        log.Info().Msg("config.yaml is not at one of these 2 locations: $HOME/.config/.shownote $HOME/.shownote")
    }
    log.Debug().Str("len of configpath", strconv.Itoa(len(configpath)))
    if len(configpath) != 0 {
        config_file,err := os.Open(configpath)
        b_config, err := io.ReadAll(config_file)
        if err != nil {
            log.Err(err).Stack().Msg("error reading bytes from config file")
        }
        err = yaml.Unmarshal(b_config, &c_config)
        if err != nil {
            log.Err(err).Stack().Msg("error unmarshalling config file")
        }
        config_file.Close()
    }

    args := os.Args

    switch len(os.Args) {
        case 2: // eg. sn tidy
            switch args[1] {
                case "tidy", "t":
                    u.Tidy(c_config.notes_folder, c_config.file_ext)
                case "-h" , "h" , "-help", "help":
                    fmt.Println("")
                    fmt.Println("usage:")
                    fmt.Println("sn 12 --> \t\topens note for task 12")
                    fmt.Println("sn add 12 --> \t\tif note for task 12 does not exist adds one, otherwise opens it")
                    fmt.Println("sn delete 12 --> \tdeletes note for task 12")
                    fmt.Println("")
                default:
                    taskno := args[1]
                    notere := regexp.MustCompile(`\d+`)
                    if notere.MatchString(taskno) {
                        u.Open_tasknote(taskno, c_config.notes_folder, c_config.file_ext)
                    } else {
                        fmt.Println("\n\033[7;31munsupported usage\033[0;0m, for help use : \"sn -h or sn help\"")
                        fmt.Println("")
                    }
            }

        case 3: // eg. sn 122 add
            method, taskno := args[1], args[2]
            switch method {
                case "add" , "ADD" , "a" , "A":
                    u.Add_tasknote(taskno, c_config.notes_folder, c_config.file_ext)
                default:
                    fmt.Println("\n\033[7;31munsupported usage\033[0;0m, for help use : \"sn -h or sn help\"")
                    fmt.Println("")
            }
        default:
            err := errors.New("No support for other than 2 or 3 arguments ")
            fmt.Println("\n\033[7;31munsupported usage\033[0;0m, for help use : \"sn -h or sn help\"")
            fmt.Println()
            log.Error().Stack().Err(err).Msg("No support for other than 2 or 3 arguments ")
            
    }
}
