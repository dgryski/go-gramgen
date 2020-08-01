//go:generate ragel-go lexer.rl
//go:generate gofmt -w lexer.go
//go:generate goyacc parser.y

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func main() {

	maxDepth := flag.Int("m", 8, "max recursion depth")
	input := flag.String("f", "", "input file")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

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

	if err := typecheck(symtab, g); err != nil {
		log.Fatal(err)
	}

	changed := true
	for changed {
		changed = false
		for k, v := range symtab {
			var b bool
			symtab[k], b = optimize(v)
			if b {
				changed = true
			}
		}
	}

	g = symtab["START"]
	cheapest(symtab, g)

	// TODO(dgryski): add range, repeat
	// TODO(dgryski): add variables to rules
	// TODO(dgryski): what else to pick from dharma syntax?
	// TODO(dgryski): common library of useful items
	// TODO(dgryski): update syntax to match cup?
	// TODO(dgryski): support "\"" and "\n" in lexer

	rand.Seed(time.Now().UnixNano())

	g.generate(os.Stdout, *maxDepth)
}

// have we visited this variable already during typecheck
var typeCache = make(map[string]error)

func typecheck(symtab map[string]generator, sym generator) error {
	// typecheck the tree rooted at sym
	// look for undefined symbols in the rules

	switch s := sym.(type) {
	case terminal:
	case intrange:
	case chrange:
	case epsilon:

	case *choice:
		for _, i := range s.c {
			if err := typecheck(symtab, i); err != nil {
				return err
			}
		}

	case *sequence:
		for _, i := range s.s {
			if err := typecheck(symtab, i); err != nil {
				return err
			}
		}

	case *variable:
		if err, ok := typeCache[s.v]; ok {
			// already recursed here
			return err
		}
		s2, ok := symtab[s.v]
		if !ok {
			return fmt.Errorf("unknown symbol: %v", s.v)
		}

		typeCache[s.v] = nil
		err := typecheck(symtab, s2)
		typeCache[s.v] = err
		return err

	default:
		panic("unknown generator type")
	}

	return nil
}

//  cache variable -> cheapest generator lookups
var seen = make(map[string]bool)
var cheapestOption = make(map[string]generator)

func cheapest(symtab map[string]generator, sym generator) (g generator, d int) {
	// typecheck the tree rooted at sym
	// look for undefined symbols in the rules

	switch s := sym.(type) {
	case terminal:
	case intrange:
	case chrange:
	case epsilon:

	case *choice:
		g, d := cheapest(symtab, s.c[0])
		for _, c := range s.c[1:] {
			if _, dd := cheapest(symtab, c); dd < d {
				g, d = c, dd
			}
		}
		s.cheap = g
		return g, d

	case *sequence:
		_, d := cheapest(symtab, s.s[0])
		for _, c := range s.s[1:] {
			if _, dd := cheapest(symtab, c); dd > d {
				d = dd
			}
		}
		return s, d + 1

	case *variable:
		if _, ok := seen[s.v]; ok {
			return sym, math.MaxUint32
		}

		ss := symtab[s.v]
		seen[s.v] = true
		g, d := cheapest(symtab, ss)
		delete(seen, s.v)
		cheapestOption[s.v] = g
		return g, d + 1

	default:
		panic("unknown generator type")
	}

	return sym, 0
}

func optimize(sym generator) (generator, bool) {
	// typecheck the tree rooted at sym
	// look for undefined symbols in the rules

	switch s := sym.(type) {
	case terminal:
	case intrange:
	case chrange:
	case epsilon:
	case *variable:

	case *choice:
		if len(s.c) == 1 {
			return optimize(s.c[0])
		}

		var changed bool
		for i, c := range s.c {
			var b bool
			s.c[i], b = optimize(c)
			changed = changed || b
		}
		return sym, changed

	case *sequence:
		if len(s.s) == 1 {
			return optimize(s.s[0])
		}

		var changed bool
		for i, c := range s.s {
			var b bool
			s.s[i], b = optimize(c)
			changed = changed || b
		}
		return sym, changed

	default:
		panic("unknown generator type")
	}

	return sym, false
}
