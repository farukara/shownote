package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
    taskno := os.Args[1]

    app := "task"
    arg0 := "_get"
    arg1 := taskno + ".uuid"

    cmd := exec.Command(app, arg0, arg1)
    task_uuid, err := cmd.Output()

    app = "nvim"
    arg0 = "~/.tasknotes/" + string(task_uuid) + ".md"

    homename, err := os.UserHomeDir()
    if err != nil {
        fmt.Println( err )
    }

    fmt.Println("task uuid:", string(task_uuid))
    arg1 = homename + "/.tasknotes/" + strings.TrimSpace(string(task_uuid)) + ".md"
    cmd = exec.Command(app, arg1)
    // stdout, err := cmd.Output()
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    err = cmd.Run()
    if err != nil {
        fmt.Println( err )
    }

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // Print the output
    // fmt.Println(string(stdout))
}

