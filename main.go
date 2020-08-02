//go:generate ragel-go lexer.rl
//go:generate gofmt -w lexer.go
//go:generate goyacc parser.y

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
)

func main() {

	maxDepth := flag.Int("m", 8, "max recursion depth")
	input := flag.String("f", "", "input file")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	benchmark := flag.Bool("bench", false, "run benchmark")
	seed := flag.Uint64("seed", 0, "xorm random seed")
	dump := flag.Bool("dump", false, "dump parse tree")
	items := flag.Int("n", 1, "number of runs")
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

	for k := range symtab {
		idx := len(symtabToIdx)
		symtabToIdx[k] = idx
	}
	symtabIdx = make([]generator, len(symtabToIdx))

	if err := typecheck(symtab, g); err != nil {
		log.Fatal(err)
	}

	changed := true
	for changed {
		changed = false
		for k, v := range symtab {
			var b bool
			symtab[k], b = optimize(v)
			changed = changed || b
		}
	}

	for k, v := range vars {
		idx := symtabToIdx[k]
		v.idx = idx
		ss := symtab[v.v]
		symtabIdx[idx] = ss
	}

	seen = make(map[string]bool)
	seen["START"] = true
	g = symtab["START"]
	unused(symtab, g)

	var remove []string
	for k := range symtab {
		if !seen[k] {
			remove = append(remove, k)
		}
	}

	for _, k := range remove {
		delete(symtab, k)
	}

	g = symtab["START"]
	seen = make(map[string]bool)
	cheapestOption = make([]generator, len(symtabIdx))
	cheapest(symtab, g)

	// TODO(dgryski): add range, repeat
	// TODO(dgryski): add variables to rules
	// TODO(dgryski): what else to pick from dharma syntax?
	// TODO(dgryski): common library of useful items
	// TODO(dgryski): update syntax to match cup?
	// TODO(dgryski): support "\"" and "\n" in lexer

	if *dump {
		var keys []string
		for k := range symtab {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Printf("%v := %v\n", k, symtab[k].String())
		}
		return
	}

	if *seed == 0 {
		*seed = uint64(time.Now().UnixNano())
	}
	xrand = xorm(*seed)

	if *benchmark {
		var buf bytes.Buffer
		var total int
		t0 := time.Now()
		for i := 0; i < 5000000; i++ {
			if i&0xfffff == 0 {
				t1 := time.Since(t0)
				log.Printf("%d bytes in %v (%s bytes/sec)", total, time.Since(t0), humanize.Bytes(uint64(float64(total)/float64(t1.Seconds()))))
				t0 = time.Now()
				total = 0
			}
			buf.Reset()
			g.generate(&buf, *maxDepth)
			total += buf.Len()
		}
		return
	}

	buf := bufio.NewWriter(os.Stdout)
	for i := 0; i < *items; i++ {
		buf.Reset(os.Stdout)
		g.generate(buf, *maxDepth)
		buf.Flush()
	}
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
var seen map[string]bool
var cheapestOption []generator

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
		cheapestOption[s.idx] = g
		return g, d + 1

	default:
		panic("unknown generator type")
	}

	return sym, 0
}

var symtabToIdx = make(map[string]int)
var symtabIdx []generator

func optimize(sym generator) (generator, bool) {
	// typecheck the tree rooted at sym
	// look for undefined symbols in the rules

	switch s := sym.(type) {
	case terminal:
	case intrange:
	case chrange:
	case epsilon:

	case *variable:
		ss := symtab[s.v]

		switch r := ss.(type) {
		case terminal:
			return r, true
		case intrange:
			return r, true
		case chrange:
			return r, true
		case epsilon:
			return r, true
		case *variable:
			return r, true
		}

		return sym, false

	case *choice:
		if len(s.c) == 1 {
			g, _ := optimize(s.c[0])
			return g, true
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
			g, _ := optimize(s.s[0])
			return g, true
		}

		var changed bool
		for i, c := range s.s {
			var b bool
			s.s[i], b = optimize(c)
			changed = changed || b
		}

		for i := 0; i < len(s.s)-1; i++ {
			t1, ok := s.s[i].(terminal)
			if !ok {
				continue
			}

			t2, ok := s.s[i+1].(terminal)
			if !ok {
				continue
			}

			s.s[i] = terminal(string(t1) + string(t2))
			s.s[i+1] = epsilon{}
			changed = true
		}

		for i := 0; i < len(s.s); i++ {
			v, ok := s.s[i].(*variable)
			if !ok {
				continue
			}

			ss := symtabIdx[v.idx]
			sseq, ok := ss.(*sequence)
			if !ok {
				continue
			}

			seq := make([]generator, 0, len(s.s)+len(sseq.s))
			seq = append(seq, s.s[:i]...)
			seq = append(seq, sseq.s...)
			seq = append(seq, s.s[i+1:]...)
			s.s = seq
			changed = true
		}

		for i := 0; i < len(s.s); i++ {
			if _, ok := s.s[i].(epsilon); !ok {
				continue
			}

			s.s = append(s.s[:i], s.s[i+1:]...)

			changed = true
		}

		return sym, changed

	default:
		panic("unknown generator type")
	}

	return sym, false
}

func unused(symtab map[string]generator, sym generator) {
	// typecheck the tree rooted at sym
	// look for undefined symbols in the rules

	switch s := sym.(type) {
	case terminal:
	case intrange:
	case chrange:
	case epsilon:

	case *choice:
		for _, i := range s.c {
			unused(symtab, i)
		}

	case *sequence:
		for _, i := range s.s {
			unused(symtab, i)
		}

	case *variable:
		if seen[s.v] {
			return
		}

		seen[s.v] = true

		s2 := symtab[s.v]
		unused(symtab, s2)

	default:
		panic("unknown generator type")
	}
}
