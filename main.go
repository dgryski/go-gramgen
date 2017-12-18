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

	rand.Seed(time.Now().UnixNano())

	g.generate(os.Stdout)
}
