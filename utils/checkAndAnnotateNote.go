package utils

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
)

// This function checks the taskNo whether it has "Notes" annotation and add one if it doesn't.
func checkAndAnnotateNote(taskNo string) error {
    // First check if "virtual tag ANNOTATED" exist on tag for early exit strategy.
    // https://taskwarrior.org/docs/dom/
    app := "task"
    cmd := exec.Command(app, taskNo + ".tags.ANNOTATED" )
    result, err := cmd.Output()
    if err != nil {
        log.Err(err).Stack().Msg("error running first cmd in function" +err.Error())
    }

    if len(result) == 0 {
        fmt.Println(taskNo, "has no annotatations, creating a \"Notes\" annotation")
        cmd := exec.Command(app, taskNo, "annotate", "Notes" )
        fmt.Println("this commad will be run:", cmd.String())
        result, err = cmd.Output()
        if err != nil {
            fmt.Println(taskNo, "error creating Notes annotation for the note:", err)
            fmt.Println("result:", string(result))
            return err
        }
        fmt.Println("task", taskNo, " is annotated...")
        return nil
    }

    // case where annotations exist. Here we are gonna check 5 notes to see of
    // any of them is Notes. We checking 5 because tw api does not return
    // number of annotations.
    for i:=1; i<5; i++ {
        // t _get 1121.annotations.1.description
        arg0 := fmt.Sprintf("%s.annotations.%d.description", taskNo, i)
        cmd := exec.Command(app, "_get", arg0)
        result,err := cmd.Output()
        if err != nil {
            fmt.Println("error running cmd during annotations check:", err)
            return err
        }
        // fmt.Println("result:", string(result))
        if string(bytes.TrimSpace(result)) == "Notes" {
            fmt.Println(taskNo, "has Notes annotation...")
            return nil
        }
    }

    // case where there are some annotations but none of them is Notes.
    // create annotation.
    cmd = exec.Command(app, taskNo, "annotate", "Notes" )
    fmt.Println("this commad will be run:", cmd.String())
    result, err = cmd.Output()
    if err != nil {
        fmt.Println(taskNo, "error creating Notes annotation for the note:", err)
        fmt.Println("result:", string(result))
        return err
    }
    fmt.Println("task", taskNo, " is annotated...")
    return nil
}
