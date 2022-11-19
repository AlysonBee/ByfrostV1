package main

import "fmt"

// C lexer body
type Token struct {
	name     string
	filepath string
	typ      string
	line     int
	spaces   int
	tabs     int

	// For function call tokens only
	label string
	id    string

	// CSS styling
	styleTag string
}

var TokenOptions = map[byte]map[byte]string{
	'<': {
		'<': "BITSHIFTLEFT",
		'=': "LESSEQU",
	},
	'>': {
		'>': "BITSHIFTRIGHT",
		'=': "GREATEREQU",
	},
	'&': {
		'&': "AND",
		'=': "ANDEQU",
	},
	'|': {
		'|': "OR",
		'=': "OREQU",
	},
	'-': {
		'>': "PTRDEREF",
		'=': "LESSEQU",
		'-': "MINUSMINUS",
	},
	'+': {
		'+': "PLUSPLUS",
		'=': "PLUSEQU",
	},
	'*': {
		'=': "MULTIASSIGN",
	},
}

// For the tokens made up of 3 or more strings.
var tokenOptionStrings = map[byte]string{
	'.': "...",
	'<': "<<=",
	'>': ">>=",
}

var reservedTokens = map[string]string{
	// Reserved
	"for":      "FOR",
	"do":       "DO",
	"while":    "WHILE",
	"if":       "IF",
	"else":     "ELSE",
	"switch":   "SWITCH",
	"case":     "CASE",
	"continue": "CONTINUE",
	"return":   "RETURN",
	"sizeof":   "SIZEOF",
	"typedef":  "TYPEDEF",

	// Datatypes
	"int":      "DATATYPE",
	"long":     "DATATYPE",
	"char":     "DATATYPE",
	"unsigned": "DATATYPE",
	"short":    "DATATYPE",
	"struct":   "STRUCT",
	"union":    "UNION",
}

var singleTokenList = map[byte]string{
	'(': "OPENBRACKET",
	')': "CLOSEBRACKET",
	'{': "OPENCURLY",
	'}': "CLOSECURLY",
	'[': "OPENBLOCK",
	']': "CLOSEBLOCK",
	';': "SEMICOLON",
	':': "COLON",
	'?': "QMARK",
	'.': "DOT",
	',': "COMMA",
	'*': "STAR",
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

type Parser struct {
	tokens []*Token
	length int
	ptr    int
	save   int
}

func (p *Parser) next() *Token {
	if p.ptr < p.length {
		t := p.tokens[p.ptr]
		p.ptr++
		return t
	}
	return nil
}

func (p *Parser) current() *Token {
	return p.tokens[p.ptr]
}

// DEBUG
func (p *Parser) whereAreWe() {
	fmt.Printf("%d of %d\n", p.ptr, p.length)
}

func (p *Parser) prev() *Token {
	if p.ptr > 0 {
		return p.tokens[p.ptr-1]
	}
	return nil
}

func (p *Parser) snapshot() {
	p.save = p.ptr
}

func (p *Parser) rollback() {
	p.ptr = p.save
}

func (p *Parser) eof() bool {
	return p.ptr >= p.length-1
}

func (p *Parser) lookAhead() *Token {
	if p.ptr+1 >= p.length {
		return nil
	}

	return p.tokens[p.ptr+1]
}

func (p *Parser) moveBy(steps int) {
	for steps > 0 {
		p.next()
		steps--
	}
}

func newParser(tokens []*Token) *Parser {
	return &Parser{
		tokens: tokens,
		length: len(tokens),
		ptr:    0,
		save:   0,
	}
}
