# Show Note for Taskwarrior

## Overview
"shownote" (or sn for short) add note-taking functionality to Taskwarrior. It requires on Taskwarrior to be installed first. Taskwarrior is a open source task management tool for terminals. You can download it [here](https://taskwarrior.org/).

Taskwarrior provides an annotation feature where you can write small notes. However it's not suitable for heavy note-taking. Enter "shownote". With show note you can add view delete a notes pertaining to a task on a totally different text file. Each task will have its own note file. 

## Installation
### Build from source you need [Go](https://go.dev/).
- download this repo by:

    $ git clone https://github.com/farukara/shownote
    cd show note
    make

make runs go build and copies the executable binary "sn" to usr/local/bin

### Directly download executable (to be implemented)
- download the file from release page 
- save file somewhere in the $PATH
- you can see folders that are on the $PATH by:

    printenv $PATH

on the console.

## Commands

- add or a
- delete or del or d (not implemented yet)
- tidy or t (not implemented yet)

## Usage
    sn 12
opens note for task number 12

    sn add 16 
adds a note for task 16. if file exist with same uuid, it just opens it

    sn 13
if you try to open a non-existent note. it will give options or  adding one and cancelling.

## Tips

- you can add following line to your bash/zsh rc file to make calling the gotask with less keys strokes: 

    alias shownote sn
for shownote. 

## Configuration
**shownote** will look for a config.yaml file at two places: `$HOME/.config/shownote/config.yaml` or `$HOME/.shownote/config.yaml`. If there is no config file, then it uses the following defaults: 

    file_ext: ".md" #dont forget initial dot 
    editor: "vim" #empty for default system $EDITOR
    notes_folder: ".tasknotes"  #relative to user home folder
    open_tasknote_after_creation: 1

You can find a sample config file in config.yaml file above. All options have to be present in the file for it to work. Basically, only change the right side of the colon in the config file.

## Advantages

*Data is yours*. You can use it however you want. If you want you can build a cloud solution, which will enable syncing.
*It's decoupled from Taskwarrior*. Notes are kept another folder and do not interfere with Taskwarrior at all.
*Extensible*. You can extend the functionality to suit your needs, such as spaced repetition.

## Improvements Needed

- Only testted and used in macOS. It should run alright in Linux. For other platforms it's not been used nor tested extensively.

## Concepts

## TODOs

- add delete
- add method to list orphan notes
- add a logger lib
- add a testing lib
- make path use path from std library to be compatible
- main 12 is opening the notes file even if it's not created before. gotta check either there is Notes annotation or if the file exist
- for recursive tasks notes annotion is inherited from the parent, but it does not point to an existing note. Make it to ask if user want to open parent note instead.
- set config.yaml with if else clauses to use different folders for dev and prod
- add options for log levels with nolog as well.
- separate functions to their own files, and move all code to src folder. only main.go remains at the top level
- adding new note when it does not exist annotates the main task with "notes:Notes", think about putting option for that into config
