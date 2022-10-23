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
    }
    notes_folder = path.Join(home, notes_folder)

    files, err := fs.Glob(os.DirFS(notes_folder), "*")
    if err != nil {
        log.Err(err).Stack().Msg("could not glob the notes folder")
        fmt.Println("could not glob the notes folder:", err )
    }


    ch := make(chan filecmd, 1000)
    wgout := &sync.WaitGroup{}
    wgout.Add(2)
    go func (ch chan filecmd) {
        defer wgout.Done()
        for _, file := range files {
            //if its not a uuid named task file
            if len(file) != 39 {
                continue
            }

            // look if task exist for that note
            app := "task"
            arg0 := "_get"
            arg1 := file[:8] + ".id"

            cmd := exec.Command(app, arg0, arg1)
            filecmd := filecmd{
                cmd: cmd,
                filename: file,
            }
            ch <- filecmd

        }
        close(ch)
    }(ch)

    go func (ch chan filecmd) {
        defer wgout.Done()
        wg := &sync.WaitGroup{}
        wg.Add(4)
        for i:=0; i<4; i++ {
            go func (ch chan filecmd) {
                defer wg.Done()
                for {
                    filecmd, ok := <- ch
                    if !ok {
                        break
                    }
                    id, err := filecmd.cmd.Output()
                    if err == nil && len(id) < 2 {
                        fmt.Println(path.Join(notes_folder, filecmd.filename))
                    }
                }
            }(ch)
        }
        wg.Wait()
    }(ch)
    wgout.Wait()
}
