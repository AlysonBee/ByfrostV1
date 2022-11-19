package x86

import "unicode"

var reservedTokens = map[string]string{
	// Reserved
	"rax": "RAX",
	"rbx": "RBX",
	"rcx": "RCX",
	"rdx": "RDX",
	"eax": "EAX",
	"ebx": "EBX",
	"ecx": "ECX",
	"edx": "EDX",
	"ax":  "AX",
	"bx":  "BX",
	"cx":  "CX",
	"dx":  "DX",
	"al":  "AL",
	"ah":  "AH",
	"bl":  "BL",
	"bh":  "BH",
	"cl":  "CL",
	"ch":  "CH",
	"dl":  "DL",
	"dh":  "DH",

	"rdi":    "RDI",
	"rsi":    "RSI",
	"sycall": "SYCALL",
	"int":    "INT",
	// Datatypes
}

func ASMscan(content string, filepath string) *Token {
	space, tab := handleWhitespace(content)

	// This was the only thing needed to be done
	if FilePtr >= len(content) {
		return &Token{}
	}

	if content[FilePtr] == ';' {
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
