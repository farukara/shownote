package utils

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/rs/zerolog/log"
)

type filecmd struct {
    cmd *exec.Cmd
    filename string
}

// find orphan notes that does not have a corresponding task
// in taskwarrior
func Tidy(notes_folder, file_ext string) {
    home,err := os.UserHomeDir()
    if err != nil {
        log.Err(err).Stack().Msg("could not get home dir")
        fmt.Println("could not get home dir:", err )
        return
    }
    notes_folder = path.Join(home, notes_folder)

    files, err := fs.Glob(os.DirFS(notes_folder), "*")
    if err != nil {
        log.Err(err).Stack().Msg("could not glob the notes folder")
        fmt.Println("could not glob the notes folder:", err )
        return
    }

    filesChan := make(chan string)
    go func() {
        for _, file := range files {
            filesChan <- file
        }
        close(filesChan)
    }()
    

    cmdWorkers := 30

    fileCmdChan := make(chan filecmd)
    // wgout := &sync.WaitGroup{}
    // wgout.Add(1)
    cmdwg := &sync.WaitGroup{}
    cmdwg.Add(cmdWorkers)
    for i:=0; i<cmdWorkers; i++ {
        go func (ch chan<- filecmd) {
            defer cmdwg.Done()
            for file := range filesChan {
                //if its not a uuid named task file
                if len(file) != 39 {
                    continue
                }

                // look if task exist for that note
                app := "task"
                arg0 := "_get"
                arg1 := file[:8] + ".uuid"

                cmd := exec.Command(app, arg0, arg1)
                filecmd := filecmd{
                    cmd: cmd,
                    filename: file,
                }
                // fmt.Println(" sending to chan")
                ch <- filecmd

            }
        }(fileCmdChan)
    }

    go func () {
        cmdwg.Wait()
        close(fileCmdChan)
    }()

    var orphanFilesChan = make(chan string)

    workers := 50
    wg := &sync.WaitGroup{}
    wg.Add(workers)

    for i:=0; i<workers; i++ {
        go func () {
            defer wg.Done()
            for filecmd := range fileCmdChan {
                // fmt.Println("received on chan............")
                // filecmd, ok := <- ch
                // if !ok {
                //     break
                // }
                uuid, err := filecmd.cmd.Output()
                if err != nil {
                    fmt.Println("\033[7;31merror occured during following command:")
                    fmt.Println(uuid, filecmd)
                    continue
                }
                if len(uuid) != 37 { // length of uuid, plus 1 for newline, not including the .md file extention
                    orphanFilesChan <- filecmd.filename
                }
            }

        }()
    }
    var orphanFiles []string

    go func() {
        wg.Wait()
        close(orphanFilesChan)
    }()

    lastwg := &sync.WaitGroup{}
    lastwg.Add(1)
    go func () {
        defer lastwg.Done()
        for orphanfile := range orphanFilesChan {
            orphanFiles = append(orphanFiles, orphanfile)
        }
    }()
    lastwg.Wait()

    if len(orphanFiles) == 0 {
        fmt.Println("No orphan notes have been found.")
    } else {
        fmt.Println("\033[7;34mFollowing files don't have an assosiated task.")
        fmt.Println("They are called Orphan Notes")
        fmt.Println("You can delete them with the following command:")
        fmt.Println("\033[7;0m")
        fmt.Println("rm  <filename>")
        fmt.Println()
        fmt.Println("\033[7;34mOrphan files:\033[7;0m")
        for i,orphanFile := range orphanFiles {
            fmt.Println(i+1, ". ", path.Join(notes_folder, orphanFile))

        }

    }
}
