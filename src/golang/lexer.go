package golang

import (
	"strings"
	"unicode"
)

var LineNumber = 1
var FilePtr = 0

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
	':': {
		'=': "DECLASSIGN",
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
	"if":       "IF",
	"else":     "ELSE",
	"switch":   "SWITCH",
	"case":     "CASE",
	"continue": "CONTINUE",
	"return":   "RETURN",
	"func":     "FUNC",
	"range":    "RANGE",
	"true":     "TRUE",
	"false":    "FALSE",
	"nil":      "NIL",

	// Datatypes
	"int":       "DATATYPE",
	"uint64":    "DATATYPE",
	"uint":      "DATATYPE",
	"string":    "DATATYPE",
	"float":     "DATATYPE",
	"struct":    "STRUCT",
	"interface": "INTERFACE",
	"map":       "MAP",
	"var":       "VAR",
	"package":   "PACKAGE",
	"import":    "IMPORT",
	"type":      "TYPE",
	"const":     "CONST",
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

type Token struct {
	Name   string
	Typ    string
	Line   int
	Spaces int
	Tabs   int

	// For function call tokens only
	Label string
	Id    string

	SymbolTab *SymbolTable

	// CSS styling
	StyleTag string

	// If it has anything, this will determine if this token is part of a
	// function call or library or whatever
	Attribute string
	Misc      string
}

func (t *Token) labelAndType(value string) {
	t.Label = value
	t.Typ = value
}

func handleWhitespace(content string) (int, int) {
	var (
		space int
		tab   int
	)

	for FilePtr < len(content) {
		switch content[FilePtr] {
		case '\t':
			tab++
		case ' ':
			space++
		case '\n':
			LineNumber++
		default:
			return space, tab
		}
		FilePtr++
	}
	return 0, 0
}

func handleThreeOrMore(content string, tabs, spaces int, typeName string) *Token {
	var currTok strings.Builder

	toCompare, ok := tokenOptionStrings[content[FilePtr]]
	if !ok {
		return nil
	}

	compCounter := 0
	contCounter := FilePtr
	for contCounter < len(content) && compCounter < len(toCompare) {
		currTok.WriteByte(content[contCounter])

		if content[contCounter] != toCompare[compCounter] {
			return nil
		}
		contCounter++
		compCounter++
	}

	FilePtr += compCounter
	return &Token{
		Name:   currTok.String(),
		Typ:    typeName,
		Line:   LineNumber,
		Tabs:   tabs,
		Spaces: spaces,
	}
}

func handleToken(token byte, tabs, spaces int) *Token {
	FilePtr++
	return &Token{
		Name:   string(token),
		Typ:    singleTokenList[token],
		Line:   LineNumber,
		Tabs:   tabs,
		Spaces: spaces,
	}
}

func handleTwoToken(content string, tabs, spaces int,
	first, second byte, firstType, secondType string) *Token {

	var (
		typ     string
		currTok strings.Builder
	)

	typ = firstType
	currTok.WriteByte(first)

	FilePtr++
	if FilePtr < len(content) && content[FilePtr] == second {
		typ = secondType
		currTok.WriteByte(second)
		FilePtr++
	}

	return &Token{
		Name:   currTok.String(),
		Typ:    typ,
		Line:   LineNumber,
		Tabs:   tabs,
		Spaces: spaces,
	}
}

func handleDigit(content string, tabs, spaces int) *Token {
	var currNum strings.Builder

	for FilePtr < len(content) {

		if !unicode.IsDigit(rune(content[FilePtr])) {
			break
		}

		currNum.WriteByte(content[FilePtr])
		FilePtr++
	}
	return &Token{
		Name:   currNum.String(),
		Typ:    "NUMBER",
		Line:   LineNumber,
		Tabs:   tabs,
		Spaces: spaces,
	}
}

func handleChar(content string, tabs, spaces int) *Token {
	var currChar strings.Builder

	// Write the first ' to prevent a premature loop exit below.
	currChar.WriteByte(content[FilePtr])

	FilePtr++
	for FilePtr < len(content) {
		if content[FilePtr] == '\'' {
			currChar.WriteByte(content[FilePtr])
			FilePtr++
			break
		}

		currChar.WriteByte(content[FilePtr])
		FilePtr++
	}

	return &Token{
		Name:   currChar.String(),
		Typ:    "CHAR",
		Line:   LineNumber,
		Tabs:   tabs,
		Spaces: spaces,
	}
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func handleID(content string, tabs, spaces int) *Token {
	var (
		currID strings.Builder
		typ    string
	)
	for FilePtr < len(content) {
		if !isAlpha(rune(content[FilePtr])) && !unicode.IsDigit(rune(content[FilePtr])) {
			break
		}

		currID.WriteByte(content[FilePtr])
		FilePtr++
	}

	typ, ok := reservedTokens[currID.String()]
	if !ok {
		typ = "ID"
	}

	return &Token{
		Name:   currID.String(),
		Typ:    typ,
		Line:   LineNumber,
		Tabs:   tabs,
		Spaces: spaces,
	}
}

func handleMultipeChoice(content string, tabs, spaces int, typeName string) *Token {
	var (
		currTok strings.Builder
		typ     string
	)

	currTok.WriteByte(content[FilePtr])
	typ = typeName

	first := content[FilePtr]
	FilePtr++
	if FilePtr < len(content) {
		opts, ok := TokenOptions[first]
		if ok {
			// TODO: Make all of this a separate function entirely. Too many nests.
			typ, ok = opts[content[FilePtr]]
			if ok {
				currTok.WriteByte(content[FilePtr])
				FilePtr++
			}
		}
	}

	return &Token{
		Name:   currTok.String(),
		Typ:    typ,
		Line:   LineNumber,
		Tabs:   tabs,
		Spaces: spaces,
	}
}

func handleString(content string, tabs, spaces int) *Token {

	var currString strings.Builder

	currString.WriteByte(content[FilePtr])

	FilePtr++
	for FilePtr < len(content) {
		if content[FilePtr] == '"' {
			currString.WriteByte(content[FilePtr])
			FilePtr++
			break
		}

		// Handling escaped double quotes.
		if content[FilePtr] == '\\' {
			currString.WriteByte(content[FilePtr])
			FilePtr++
			// Move one up and ignore the check.
		}

		currString.WriteByte(content[FilePtr])
		FilePtr++
	}

	return &Token{
		Name:   currString.String(),
		Typ:    "STRING",
		Line:   LineNumber,
		Tabs:   tabs,
		Spaces: spaces,
	}
}

func handleComments(content string, tabs, spaces int) *Token {
	var comment strings.Builder

	for FilePtr < len(content) {
		if content[FilePtr] == '\n' {
			LineNumber++
			break
		}
		comment.WriteString(string(content[FilePtr]))
		FilePtr++
	}

	return &Token{
		Name:   comment.String(),
		Typ:    "COMMENT",
		Line:   LineNumber,
		Tabs:   tabs,
		Spaces: spaces,
	}
}

func GOscan(content string, filepath string) *Token {
	space, tab := handleWhitespace(content)

	// This was the only thing needed to be done
	if FilePtr >= len(content) {
		return &Token{}
	}

	if content[FilePtr] == '/' {
		if FilePtr+1 < len(content) && content[FilePtr+1] == '/' {
			return handleComments(content, tab, space)
		}
	}

	if unicode.IsDigit(rune(content[FilePtr])) {
		return handleDigit(content, tab, space)
	}

	if isAlpha(rune(content[FilePtr])) {
		return handleID(content, tab, space)
	}

	switch content[FilePtr] {
	case '"':
		return handleString(content, tab, space)
	case '\'':
		return handleChar(content, tab, space)
	case '=':
		return handleTwoToken(content, tab, space, '=', '=', "ASSIGN", "COMPARE")
	case '>':
		// Make into a function
		token := handleThreeOrMore(content, tab, space, "BITSHIFTASSIGN")
		if token == nil {
			return handleMultipeChoice(content, tab, space, "GREATER")
		}
		return token
	case '<':
		token := handleThreeOrMore(content, tab, space, "BITSHIFTASSIGN")
		if token == nil {
			return handleMultipeChoice(content, tab, space, "LESSTHAN")
		}
		return token
	case '!':
		return handleTwoToken(content, tab, space, '!', '=', "NOT", "NOTEQU")
	case '+':
		return handleMultipeChoice(content, tab, space, "PLUS")
	case '-':
		return handleMultipeChoice(content, tab, space, "MINUS")
	case '%':
		return handleTwoToken(content, tab, space, '%', '=', "MODULO", "MODULOEQU")
	case '*':
		return handleTwoToken(content, tab, space, '*', '=', "STAR", "STAREQU")
	case '/':
		return handleTwoToken(content, tab, space, '/', '=', "DIV", "DIVEQU")
	case ':':
		return handleTwoToken(content, tab, space, ':', '=', "COLON", "EQUASSIGN")
	case '&':
		return handleMultipeChoice(content, tab, space, "BITAND")
	case '|':
		return handleMultipeChoice(content, tab, space, "BITOR")

	case '.':
		// Make into a function
		token := handleThreeOrMore(content, tab, space, "ELLIPSE")
		if token == nil {
			return handleToken(content[FilePtr], tab, space)
		}
		return token
	}

	return handleToken(content[FilePtr], tab, space)
}
