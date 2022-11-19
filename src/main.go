package main

import (
	"byfrostV1/errors"
	"byfrostV1/filesys"
	"byfrostV1/golang"
	"byfrostV1/server"
	"byfrostV1/styling"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	entrypoint          = flag.String("entry", "main", "the function to start scanning from")
	absoluteProjectPath = flag.String("project", "", "the go project path (used to determine imports)")
	toScan              = flag.String("dir", ".", "the directory to start scanning from")
)

func help() {
	fmt.Println("Usage: ")
	flag.PrintDefaults()
}

var GlobalNamespace string

type ImportNamespace struct {
	Path   string
	FTable map[string]*FunctionTable
	Types  map[string]*TypeList
	Files  map[string]bool
}

var Namespace = map[string]ImportNamespace{}

func GolangTokenToToken(tkn *golang.Token) *Token {
	return &Token{
		name:   tkn.Name,
		typ:    tkn.Typ,
		line:   tkn.Line,
		spaces: tkn.Spaces,
		tabs:   tkn.Tabs,
	}
}

func GenericTokensToGolangTokens(genTokens []*Token) []*golang.Token {
	var tokens []*golang.Token

	for _, tkn := range genTokens {
		tokens = append(tokens, &golang.Token{
			Name:   tkn.name,
			Typ:    tkn.typ,
			Line:   tkn.line,
			Spaces: tkn.spaces,
			Tabs:   tkn.tabs,
		})
	}
	return tokens
}

func GolangTokenToGenericTokens(golangTokens []*golang.Token) []*Token {
	var tokens []*Token

	for _, tkn := range golangTokens {
		tokens = append(tokens, &Token{
			name:   tkn.Name,
			typ:    tkn.Typ,
			label:  tkn.Label,
			line:   tkn.Line,
			spaces: tkn.Spaces,
			tabs:   tkn.Tabs,
		})
	}
	return tokens
}

func organizeData(filename string) {
	fileContent := readFile(filename)
	// BuildTokenList := supportedLanguages[targetLanguage]

	tokens := BuildTokenList(fileContent, filename)

	golang.ParseTypes(GenericTokensToGolangTokens(tokens), filename)
	golang.Parser(GenericTokensToGolangTokens(tokens), filename, *absoluteProjectPath)
}

func createPath(path string) string {
	x := 0
	for i, c := range path {
		if c == '/' {
			x = i
			break
		}
	}

	return *absoluteProjectPath + path[x:]
}

func determineLastRoot(fullPath string, function string) string {
	segments := strings.Split(fullPath, "/")
	index := len(segments) - 1
	for index > -1 {

		lastWord := strings.Split(segments[index], ".")
		if len(lastWord) > 1 {
			return lastWord[0] + "." + function
		} else {
			index--
		}
	}
	return function
}

func SearchImport(packageName string, function string, fullPath string) (*server.BodyJSON, errors.ReturnCode) {
	var absPath string

	testFunction := strings.Split(function, ".")
	if len(testFunction) == 1 {
		function = determineLastRoot(fullPath, function)
	}

	if packageName == *toScan {
		GlobalNamespace = *toScan
		FTable = Namespace[*toScan].FTable
		TTable = Namespace[*toScan].Types
		x := strings.Split(function, ".")
		if len(x) > 1 {
			packageName = x[0]
		} else {
			return nil, errors.ReturnCode(0)
		}
	}

	for _, imp := range golang.Imports {
		if imp == nil {
			continue
		}
		if packageName == imp.Alias || packageName == imp.BaseName {
			absPath = createPath(imp.Path[1 : len(imp.Path)-1])
			break
		}
	}

	if len(absPath) == 0 {
		return nil, errors.ReturnCode(0)
	}

	namespace, ok := Namespace[absPath]
	if ok {
		GlobalNamespace = absPath
		FTable = namespace.FTable
		TTable = namespace.Types
		return nil, errors.ReturnCode(0)
	}

	filesys.Files = make(map[string]bool)
	filesys.RecurseProject(absPath, []string{".go"})

	for file := range filesys.Files {
		organizeData(file)
	}

	GlobalNamespace = absPath
	Namespace[absPath] = ImportNamespace{
		Path:   absPath,
		Files:  filesys.Files,
		Types:  golangTypeListToGenericTypeList(),
		FTable: golangFTableToGenericFTable(),
	}
	TTable = Namespace[absPath].Types
	FTable = Namespace[absPath].FTable
	return nil, errors.ReturnCode(0)
}

func main() {
	flag.Parse()
	if len(*toScan) == 0 || len(*entrypoint) == 0 || len(*absoluteProjectPath) == 0 {
		help()
		os.Exit(1)
	}

	filesys.RecurseProject(*toScan, []string{".go"})

	for file := range filesys.Files {
		organizeData(file)
	}

	styling.PrepSyntaxHighlighting("styling/configs/go.json")

	GlobalNamespace = *toScan

	Namespace[*toScan] = ImportNamespace{
		Path:   *toScan,
		Files:  filesys.Files,
		Types:  golangTypeListToGenericTypeList(),
		FTable: golangFTableToGenericFTable(),
	}
	FTable = Namespace[*toScan].FTable
	TTable = Namespace[*toScan].Types
	CodeServer()
}
