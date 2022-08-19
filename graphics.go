package main

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
func Graphics() func() string {
	i := 0

	return func() string {
		nextGraphic := graphics[i]
		i++
		// fmt.Println(nextGraphic)
		return nextGraphic
	}
}
