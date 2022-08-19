package main

import "errors"

/*
	Provide some basic graphics to use
*/

var graphics = [8]string{`
    ╭───────╮
    │       │
    │
    │
    │
    │
    `,
	`
    ╭───────╮
    │       │
    │       ◯
    │
    │
    │
	`,
	`
    ╭───────╮
    │       │
    │       ◯
    │       │
    │
    │
	`,
	`
    ╭───────╮
    │       │
    │       ◯
    │       │╲
    │
    │
	`,
	`
    ╭───────╮
    │       │
    │       ◯
    │      ╱│╲
    │
    │
	`,
	`
    ╭───────╮
    │       │
    │       ◯
    │      ╱│╲
    │       │
    │
	`,
	`
    ╭───────╮
    │       │
    │       ◯
    │      ╱│╲
    │       │
    │        ╲
	`,
	`
    ╭───────╮
    │       │
    │       ◯
    │      ╱│╲
    │       │
    │      ╱ ╲
	`,
}

// Closure alert!
// Every time this is called, the next graphic from above is returned
// Basically a Python generator...
func Graphics() func() (string, error) {
	i := 0
	cap := len(graphics)

	return func() (string, error) {
		if i == cap {
			return "", errors.New("")
		} else {
			nextGraphic := graphics[i]
			i++
			return nextGraphic, nil
		}
	}
}
