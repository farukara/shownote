package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func get_task_uuid(taskno string) ([]byte, error){
    app := "task"
    arg0 := "_get"
    arg1 := taskno + ".uuid"

    cmd := exec.Command(app, arg0, arg1)
    // NOTE: make sure what is returned by err from cmd
    task_uuid, err := cmd.Output()
    return task_uuid, err
}

func open_tasknote (taskno string) {
    task_uuid, err := get_task_uuid(taskno)
    if err != nil {
        log.Println( err )
    }

    app := "nvim"
    homename, err := os.UserHomeDir()
    if err != nil {
        log.Println( err )
    }
    log.Println("task uuid:", string(task_uuid))
    arg0 := homename + "/.tasknotes/" + strings.TrimSpace(string(task_uuid)) + ".md"
    cmd := exec.Command(app, arg0)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    err = cmd.Run()
    if err != nil {
        log.Println( err )
    }
}

func add_tasknote(taskno, method string) {
    log.Println("taskno:", taskno)
    log.Println("method:", method)
    task_uuid, err := get_task_uuid(taskno)
    if err != nil {
        log.Println( err )
    }
    if method != "add" {
        log.Println(method, "is not supported")
        log.Println("gotask help for available methods")
    }
    homename, err := os.UserHomeDir()
    if err != nil {
        log.Println( err )
    }
    log.Println("task uuid:", string(task_uuid))
    filename := homename + "/.tasknotes/" + strings.TrimSpace(string(task_uuid)) + ".md"
    // filename := homename + "/Documents/golang/tasknotes/" + strings.TrimSpace(string(task_uuid)) + ".md"
    filePath, _ := filepath.Abs(filename)
    os.Create(filePath)
    log.Println("file created:", filename)
    open_tasknote(taskno)
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
            log.Println("no support for more than 2 arguments")
            
    }
}

