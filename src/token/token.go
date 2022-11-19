package token

import "fmt"

type Token struct {
	name     string
	filepath string
	typ      string
	line     int
	spaces   int
	tabs     int

	// For function call tokens only
	label string
}

func (t *Token) String() {
	fmt.Printf(
		"name   : %s\n"+
			"file   : %s\n"+
			"type   : %s\n"+
			"line   : %d\n"+
			"spaces : %d\n"+
			"tabs   : %d\n",
		t.name, t.filepath, t.typ,
		t.line, t.spaces,
		t.tabs,
	)
}

func newToken()
