package utils

import (
	"errors"
	"os/exec"
)

func Get_task_uuid(taskno string) ([]byte, error){
    app := "task"
    arg0 := "_get"
    arg1 := taskno + ".uuid"

    cmd := exec.Command(app, arg0, arg1)
    // NOTE: find out what is returned by err from cmd
    task_uuid, err := cmd.Output()
    if err != nil {
        return task_uuid, err
    }
    // returns line feed  single character "byte(10)" for non-existing tasks
    if len(task_uuid) < 2 {
        err = errors.New("no task with that number")
    }
    return task_uuid, err
}
