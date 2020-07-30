//go:generate ragel -Z lexer.rl
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
	tokens := lex(data)

	l := fuzzLexer(tokens)

	if yyParse(&l) != 0 {
		fmt.Println("parse error")
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

	rand.Seed(time.Now().UnixNano())

	g.generate(os.Stdout)
}
