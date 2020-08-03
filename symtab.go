package main

type symbolTable struct {
	rules map[string]generator
	vars  map[string]*variable // listing of all variables name -> *ptr
}

func newSymbolTable() *symbolTable {
	return &symbolTable{
		rules: map[string]generator{"EMPTY": epsilon{}},
		vars:  make(map[string]*variable),
	}
}
