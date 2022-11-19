package main

import "fmt"

func repeatingChar(c string, times int) {
	counter := 0

	for counter < times {
		fmt.Printf(c)
		counter++
	}
}

func newline(linenumber int) {
	fmt.Printf("\n%d ", linenumber)
}

func visualize(tokens []*Token) {
	currLine := 1

	for _, tk := range tokens {
		if tk.line > currLine {
			newline(tk.line)
			currLine = tk.line
		}
		if tk.spaces > 0 {
			repeatingChar(" ", tk.spaces)
		}
		if tk.tabs > 0 {
			repeatingChar("\t", tk.tabs)
		}
		if tk.label == "DECL" {
			fmt.Println("DECL")
		}
		if tk.label == "STRUCT" {
			fmt.Println("STRUCT")
		}
		fmt.Printf(tk.name)
	}
}
