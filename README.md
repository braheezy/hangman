# Hangman
This is Go TUI [Hangman](https://www.wikihow.com/Play-Hangman) game built with the lovely [BubbleTea](https://github.com/charmbracelet/bubbletea) framework.

This project exists to teach myself Go and learn about BubbleTea because it looks like an understandable, powerful, and good looking TUI framework. My past attempt at a [TUI in Python](https://github.com/braheezy/pyrdle) didn't go so well. It's already going much better with this project :)

## Usage
For now, clone the project and then run:
```console
go run .
```

Enjoy!

## Status
The following is to be implemented:
- [x] End game lose condition :heavy_multiplication_x:
    - Yes, right now if you can't solve the word, the program crashes
- [ ] Show guessed letters :b::a:
    - The user can't see which letters are guessed which is awkward when playing
- [ ] Beautification :sunglasses:
    - The Board is the only thing with style right now
    - Add Style to:
        - [ ] Hangman graphic
        - [ ] User input area
        - [ ] Header/Footer
        - [ ] Game message area. Distinct warnings vs success message styles
- [x] Clear terminal screen :boom:
    - Before launching, clear the entire screen for maximum cleanliness
- [x] Sanitize better :earth_americas:
    - Characters like `.` and nothing are deemed okay. That's stupid