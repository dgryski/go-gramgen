//go:generate ragel-go lexer.rl
//go:generate gofmt -w lexer.go
//go:generate goyacc parser.y

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {

	input := flag.String("f", "", "input file")
	flag.Parse()

	var data []byte
	var err error
	if *input != "" {
		data, err = ioutil.ReadFile(*input)
	} else {
		data, err = ioutil.ReadAll(os.Stdin)
	}

	if err != nil {
		log.Fatalf("error reading grammar: %v", err)
	}

	l := lex(data)

	if yyParse(l) != 0 {
		// don't log; yyError() has already done so
		return
	}

	g, ok := symtab["START"]
	if !ok {
		log.Fatal("unable to find START")
	}

	err := typecheck(symtab, g)
	if err != nil {
		log.Fatal(err)
	}

	// TODO(dgryski): add range, repeat
	// TODO(dgryski): add variables to rules
	// TODO(dgryski): what else to pick from dharma syntax?
	// TODO(dgryski): common library of useful items
	// TODO(dgryski): add optimization pass to remove extra nonterminal -> terminal steps
	// TODO(dgryski): add maxdepth param
	// TODO(dgryski): add "cheapest non-terminal" table for when max depth is exceeded
	// TODO(dgryski): add better error messages for parsing ruleset
	// TODO(dgryski): update syntax to match cup?
	// TODO(dgryski): support "\"" and "\n" in lexer

	rand.Seed(time.Now().UnixNano())

	g.generate(os.Stdout)
}

// have we recursed on this variable already
var typeCache = make(map[string]error)

func typecheck(symtab map[string]generator, sym generator) error {
	// typecheck the tree rooted at sym
	// look for undefined symbols in the rules

	switch s := sym.(type) {
	case terminal:
	case intrange:
	case chrange:
	case epsilon:

	case choice:
		for _, i := range s {
			if err := typecheck(symtab, i); err != nil {
				return err
			}
		}

	case sequence:
		for _, i := range s {
			if err := typecheck(symtab, i); err != nil {
				return err
			}
		}

	case variable:
		if err, ok := typeCache[string(s)]; ok {
			// already recursed here
			return err
		}

		s2, ok := symtab[string(s)]
		if !ok {
			return fmt.Errorf("unknown symbol: %v", s)
		}

		typeCache[string(s)] = nil
		err := typecheck(symtab, s2)
		typeCache[string(s)] = err
		return err

	default:
		panic("unknown generator type")
	}

	return nil
}
