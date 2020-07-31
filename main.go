//go:generate ragel-go lexer.rl
//go:generate gofmt -w lexer.go
//go:generate goyacc parser.y

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {

	data, _ := ioutil.ReadAll(os.Stdin)
	l := lex(data)

	if yyParse(l) != 0 {
		// don't log; yyError() has already done so
		return
	}

	g, ok := symtab["START"]
	if !ok {
		log.Fatal("unable to find START")
	}

	// TODO(dgryski): type-check and ensure all non-terminals are defined
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
