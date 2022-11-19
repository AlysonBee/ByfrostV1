package golang

import "fmt"

type Variable struct {
	Name     string
	Datatype []*Token
	Line     int
	Ptr      uint64
	Value    string
	Tokens   []*Token
}

type SymbolTable struct {
	Scope     string
	Variables []Variable
	NextTable *SymbolTable
	PrevTable *SymbolTable
}

var SymTab = NewSymbolTable("/main") // This should be the entry value the user sets.

func (s *SymbolTable) PrintVariables() {
	for _, variables := range s.Variables {
		fmt.Printf("\tvar : %s\n", variables.Name)
	}
}

func PrintVariables(s *SymbolTable) {
	prev := s
	for prev != nil {
		for _, variables := range prev.Variables {
			fmt.Printf("\tvar : %s\n", variables.Name)
		}
		prev = prev.PrevTable
	}
}

func NewSymbolTable(scope string) *SymbolTable {
	return &SymbolTable{
		Scope:     scope,
		PrevTable: nil,
		NextTable: nil,
	}
}

func (s *SymbolTable) AddVariable(tokens []*Token) {
	s.Variables = append(s.Variables, Variable{
		Name:     tokens[0].Name,
		Datatype: tokens[1:],
		Line:     tokens[0].Line,
		Value:    "",
	})
}

func PushTable(s *SymbolTable, scope string) *SymbolTable {
	newName := s.Scope + "/" + scope
	s.NextTable = NewSymbolTable(newName)
	s.NextTable.PrevTable = s
	s = s.NextTable
	return s
}

func PopTable(s *SymbolTable) *SymbolTable {
	s = s.PrevTable
	s.NextTable = nil
	return s
}

func SnapshotTable(s *SymbolTable) *SymbolTable {
	trav := s
	var newBuild *SymbolTable

	for trav.NextTable != nil {
		newBuild = &SymbolTable{
			Scope:     trav.Scope,
			Variables: trav.Variables,
			NextTable: trav.NextTable,
			PrevTable: trav.PrevTable,
		}
		newBuild = newBuild.NextTable
		trav = trav.NextTable
	}

	return newBuild
}
