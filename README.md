# Go Notes for TaskWarrior
## Overview
## Concepts
## Commands

- add or o
- delete or del or d
- orphan or o

## Usage
    gonote 12
open note for task number 12

    gonote add 16 
add a note for task 16. if file exist with same uuid, it just opens it

## Tips

- you can add following line to your bash/zsh rc file to make calling the gotask with less keys strokes: alias gonotes to. to for task open. 

## TODOs

- add method to list orphan notes
- add a logger lib
- add a testing lib
- make path use path from std library to be compatible
- main 12 is opening the notes file even if it's not created before. gotta check either there is Notes annotation or if the file exist
- when the note does not exist error should be in red and options in blue.
- for recursive tasks notes annotion is inherited from the parent, but it does not point to an existing note. Make it to ask if user want to open parent note instead.
