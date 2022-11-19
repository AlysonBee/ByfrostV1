package main

import (
	"byfrostV1/golang"
)

// Holds a list of all the functions in the program.
type FunctionTable struct {
	name          string
	entry         bool
	file          string
	chldFunctions []string
	tokens        []*Token
}

func newFTable(name, filename string, tokens []*Token) *FunctionTable {
	entry := false

	// optino to set an entry point will come when this actually properly works.
	if filename == "main" {
		entry = true
	}

	return &FunctionTable{
		name:   name,
		entry:  entry,
		file:   filename,
		tokens: tokens,
	}
}

func pushFtable(name string, filename string, tokens []*Token) {
	FTable[name] = newFTable(name, filename, tokens)
}

// Holds a list of all the functions in the program.
var FTable = map[string]*FunctionTable{}
var TTable = map[string]*TypeList{}

type StructTable struct {
	name   string
	file   string
	tokens []*Token
}

func golangFTableToGenericFTable() map[string]*FunctionTable {
	var ftable = map[string]*FunctionTable{}

	for _, item := range golang.GoFunctionTable {
		ftable[item.Name] = newFTable(
			item.Name,
			item.File,
			GolangTokenToGenericTokens(item.Tokens),
		)
	}

	golang.GoFunctionTable = map[string]*golang.FunctionItem{}
	return ftable
}

type TypeList struct {
	name   string
	file   string
	tokens []*Token
}

func newTypelis(name string, file string, tokens []*Token) *TypeList {
	return &TypeList{
		name:   name,
		file:   file,
		tokens: tokens,
	}
}

func golangTypeListToGenericTypeList() map[string]*TypeList {
	var typelist = map[string]*TypeList{}

	for _, typing := range golang.Types {
		typelist[typing.Name] = newTypelis(
			typing.Name,
			typing.File,
			GolangTokenToGenericTokens(typing.Tokens),
		)
	}

	golang.Types = map[string]*golang.TypeList{}
	return typelist
}
