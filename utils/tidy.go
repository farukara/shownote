package utils

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"

	"github.com/rs/zerolog/log"
)

// find orphan notes that does not have a corresponding task
// in taskwarrior
func Tidy(notes_folder, file_ext string) {
    home,err := os.UserHomeDir()
    if err != nil {
        log.Err(err).Stack().Msg("could not get home dir")
        fmt.Println("could not get home dir:", err )
    }
    notes_folder = path.Join(home, notes_folder)

    files, err := fs.Glob(os.DirFS(notes_folder), "*")
    if err != nil {
        log.Err(err).Stack().Msg("could not glob the notes folder")
        fmt.Println("could not glob the notes folder:", err )
    }

    for _, file := range files {
        //if its not a task file
        if len(file) != 39 {
            continue
        }

        // look if task exist for that note
                
        app := "task"
        arg0 := "_get"
        arg1 := file[:8] + ".id"

        cmd := exec.Command(app, arg0, arg1)
        id, err := cmd.Output()
        /* fmt.Println(len(string(id)))
        fmt.Println(len(id)) */
        if err == nil && len(id) < 2 {
            fmt.Println(path.Join(notes_folder, file))
        }
    } 
}
