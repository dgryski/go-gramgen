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

func (sym *symbolTable) StartRule() generator {
	return sym.rules["START"]
}

type duplicateRuleError string

func (e duplicateRuleError) Error() string {
	return "duplicate rule: " + string(e)
}

func (sym *symbolTable) addRule(r string, g generator) error {
	if _, ok := sym.rules[r]; ok {
		return duplicateRuleError(r)
	}

	sym.rules[r] = g
	return nil
}

func (sym *symbolTable) addVariable(v string) *variable {
	vptr, ok := sym.vars[v]
	if !ok {
		vptr = &variable{v: v}
		sym.vars[v] = vptr
	}

	return vptr
}
