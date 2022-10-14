# Show Note for Taskwarrior

## Overview
"shownote" (or sn for short) add note-taking functionality to Taskwarrior. It requires on Taskwarrior to be installed first. Taskwarrior is a open source task management tool for terminals. You can download it [here](https://taskwarrior.org/).

Taskwarrior provides an annotation feature where you can write small notes. However it's not suitable for heavy note-taking. Enter "shownote". With show note you can add view delete a notes pertaining to a task on a totally different text file. Each task will have its own note file. 

## Advantages

*Data is yours*. You can use it however you want. If you want you can build a cloud solution, which will enable syncing.
*It's decoupled from Taskwarrior*. Notes are kept another folder and do not interfere with Taskwarrior at all.
*Extensible*. You can extend the functionality to suit your needs, such as spaced repetition.

## Improvements Needed

- Only testted and used in Mac OS. It needs to be tested and used in other platforms.
## Concepts
## Commands

- add or a
- delete or del or d (not implemented yet)
- orphan or o (not implemented yet)

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

## TODOs

- add delete
- add method to list orphan notes
- add a logger lib
- add a testing lib
- make path use path from std library to be compatible
- main 12 is opening the notes file even if it's not created before. gotta check either there is Notes annotation or if the file exist
- for recursive tasks notes annotion is inherited from the parent, but it does not point to an existing note. Make it to ask if user want to open parent note instead.
- set config.yaml with if else clauses to use different folders for dev and prod
