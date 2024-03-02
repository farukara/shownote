# Show Note for Taskwarrior

## Overview
"shownote" (or sn for short) adds note-taking functionality to Taskwarrior in macOS and Linux. It's not tested nor used on Windows. It requires Taskwarrior to be installed first. Taskwarrior is an open source cli task manager for terminals. You can download it [here](https://taskwarrior.org/).

Taskwarrior provides an annotation feature where you can write small notes. However it's not suitable for heavy note-taking. Enter "shownote". With show note you can add view delete a notes file and attach it to a task. Each task will have its own note file.

## Installation
### Build from source you need [Go](https://go.dev/).
- download this repo by:

    $ git clone https://github.com/farukara/shownote
    cd shownote
    make

make runs go build and copies the executable binary "sn" to usr/local/bin

### Directly download executable (not available yet)
- download the file from release page 
- save file somewhere in the $PATH
- you can see folders that are on the $PATH by:

    echo $PATH

on the console. some regular path folders are:
usr/local/bin
usr/bin
/bin

sane sn file in any of this folders. Then you can use it.

## Commands

- ktask number
- task number + add or a
- task number + delete or del or d (not implemented yet, instead use tidy)
- tidy or t

## Usage
    sn 12
opens note file for task number 12. if note file for task 12 does not exist asks if you want to create one. If you choose to create then a new notes file will be created and the task will annotated with "Notes". When you list tasks, this annotation visiually indicates that task has a notes file.

    sn add 16 
adds a note for the task 16. if file exist with same uuid, it just opens it.

    sn tidy
prints out the notes that don't have a corresponding task. This happens when a task is deleted but associatted task note is not. tidy command only looks at the files with proper uuid names. so if there are other files in the same folder you have to clean them up manually. 

if there are no orphan notes "tidy" returns nothing.

"tidy" runs slow because of extensive API calls to taskwarrior (about 10s for 1,000 notes).

"tidy" only prints the orphan file names. if you want to delete them, you can pipe them into rm with the following command:

    sn tidy | xargs rm

## Tips
## Configuration
when you make from the source, the default config file will be copied to `$HOME/.config/shownote` folder. It's important that you don't change the key (left side) of each option.

**shownote** will look for a config.yaml file at two places: `$HOME/.config/shownote/config.yaml` or `$HOME/.shownote/config.yaml`. If there is no config file, then it uses the following defaults: 

    file_ext: ".md" #dont forget initial dot 
    editor: "vim" #empty for default system $EDITOR
    notes_folder: "$HOME/.tasknotes" 
    open_tasknote_after_creation: 1

You can find a sample config file in config.yaml file above. All options have to be present in the file for it to work. Again, basically, only change the right side of the colon in the config file, if you need to modify.

## Advantages

*Data is yours*. You can use it however you want. If you want you can build a cloud solution, which will enable syncing.
*It's decoupled from Taskwarrior*. Notes are kept in another folder and do not interfere with Taskwarrior at all, except for annotating the task when adding notes.
*Extensible*. You can extend the functionality to suit your needs, such as spaced repetition.

## Improvements Needed

- Only testted and used in macOS. It should run alright in Linux. For other platforms it's not been used or tested.

## Concepts

## TODOs

- add prune: on top of tidy it deletes the orphan notes
- add a testing lib
- for recursive tasks, annotion is inherited from the parent, but it does not point to an existing note. Make it to ask if user want to open parent note instead.
- set config.yaml with if else clauses to use different folders for dev and prod
- add options for log levels with nolog as well.
- currently notes are being deleted manually by user. If delete funcion will be created annotationn should be deleted as well.
