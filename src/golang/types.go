package golang

import (
	"fmt"
	"strings"
)

type FunctionItem struct {
	Name   string
	Entry  bool
	File   string
	Tokens []*Token
}

var GoFunctionTable = map[string]*FunctionItem{}

func PushFunctionItem(name string, tokens []*Token, file string) {
	GoFunctionTable[name] = &FunctionItem{
		Name:   name,
		Tokens: tokens,
		File:   file,
	}
}

type ImportList struct {
	Path     string
	Alias    string
	BaseName string
}

var declarationKeyword = map[string]bool{
	"FUNC":  true,
	"TYPE":  true,
	"CONST": true,
	"VAR":   true,
}

func printImports() {
	for _, imp := range Imports {
		fmt.Printf("Path: %s\nAlias %s\nBaseName: %s\n",
			imp.Path,
			imp.Alias,
			imp.BaseName)
	}
}

var Imports = map[string]*ImportList{}

func addImport(path, alias, baseName string) {
	Imports[baseName] = &ImportList{
		Path:     path,
		Alias:    alias,
		BaseName: baseName,
	}
}

func extractImportList(tokens []*Token, globalPointer int) {
	var (
		alias    string
		path     string
		baseName string
	)

	for globalPointer < len(tokens) {
		tk := tokens[globalPointer]
		if tk.Typ == "ID" {
			alias = tk.Name
			globalPointer++
			if globalPointer < len(tokens) {
				path = tokens[globalPointer].Name
				if len(path) > 1 {
					baseName = path[1 : len(path)-1]
				}
			}
			addImport(path, alias, baseName)
		} else if tk.Typ == "STRING" {
			alias = ""
			path = tk.Name
			pathSegments := strings.Split(tk.Name, "/")
			if len(pathSegments) > 1 {
				baseName = pathSegments[len(pathSegments)-1][0 : len(pathSegments[len(pathSegments)-1])-1]
			} else {
				baseName = pathSegments[len(pathSegments)-1][1 : len(pathSegments[len(pathSegments)-1])-1]
			}
			addImport(path, alias, baseName)
		} else if declarationKeyword[tk.Typ] {
			break
		}
		globalPointer++
	}
}

func GetImports(tokens []*Token) {
	globalPointer := 0
	for globalPointer < len(tokens) {
		tk := tokens[globalPointer]
		if tk.Name == "import" {
			globalPointer++
			break
		}
		globalPointer++
	}
	extractImportList(tokens, globalPointer)
}

type TypeList struct {
	Name    string
	File    string
	Tokens  []*Token
	Methods map[string]*FunctionItem
}

func PushMethod(name string, functionName string, tokens []*Token) {
	_, ok := Types[name]
	if ok {
		if Types[name].Methods == nil {
			Types[name].Methods = make(map[string]*FunctionItem)
		}
		Types[name].Methods[functionName] = &FunctionItem{
			Name:   functionName,
			Tokens: tokens,
			File:   Types[name].File,
		}
	}
}

func PrintAllMethods() {
	for _, name := range Types {
		fmt.Println("object name : ", name.Name)
		fmt.Println("file : ", name.File)
		fmt.Println("Tokens ")
		for _, function := range name.Methods {
			visualize(function.Tokens)
			fmt.Println("\n==========")
		}
		fmt.Println("\n")
	}
}

var Types = map[string]*TypeList{}

func PushType(name string, file string, tokens []*Token) {
	Types[name] = &TypeList{
		Name:   name,
		File:   file,
		Tokens: tokens,
	}
}

type GlobalVariables struct {
	Name   string
	File   string
	Tokens []*Token
}

type ClassList struct {
	Name       string
	File       string
	ReturnType []*Token
	Member     map[string][]*Token
	Methods    map[string][]*FunctionItem
}

/*
type Test struct {
	variable int
	naming string
}

func (t *Test) person(name string) int {
	t.naming = name
	return 1
}

ClassList {
	Name : "Test"
	File : "None.c"
	ReturnType : "struct"
	Member : [
	 	"variable" {
			type: int
		},
		"naming": {
			type: "string"
		}
	}
	Methods: [
		FunctionItem {
			Name: "person",
			File: "None.c",
			Tokens: []
		}
	]
}

var b Test
	b -> variable:/main/b.person()
		if variable -> search in main for "b" and get its datatype
			-> b = Backends
				-> search Classlist for Backends
					-> search Members for "person"
						- if not found, search Methods for "person",
							- if found, return token list



b.person().again() -> variableDeref:main/Backends.again()


b.person("timoty")

*/
