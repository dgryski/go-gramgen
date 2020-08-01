package main

import (
	"io"
	"log"
	"math/rand"
	"strconv"
)

type generator interface {
	generate(w io.Writer, depth int)
}

type terminal string

func (t terminal) generate(w io.Writer, depth int) {
	w.Write([]byte(t))
}

type variable struct {
	v string
}

func (v variable) generate(w io.Writer, depth int) {
	if depth <= 0 {
		cheapestOption[v.v].generate(w, 0)
		return
	}
	g, ok := symtab[v.v]
	if !ok {
		log.Fatalf("unknown variable %q", v.v)
	}
	g.generate(w, depth-1)
}

type choice struct {
	c     []generator
	cheap generator
}

func (c choice) generate(w io.Writer, depth int) {
	if depth <= 0 {
		c.cheap.generate(w, depth-1)
		return
	}
	n := rand.Intn(len(c.c))
	c.c[n].generate(w, depth-1)
}

func (c *choice) add(g generator) {
	c.c = append(c.c, g)
}

type sequence struct {
	s []generator
}

func (s sequence) generate(w io.Writer, depth int) {
	for _, ss := range s.s {
		ss.generate(w, depth-1)
	}
}

func (s *sequence) add(g generator) {
	s.s = append(s.s, g)
}

type intrange struct {
	low, high int
}

func (ir intrange) generate(w io.Writer, depth int) {
	n := rand.Intn(ir.high - ir.low)
	w.Write(strconv.AppendInt(nil, int64(ir.low+n), 10))
}

type chrange struct {
	low, high int
}

func (ch chrange) generate(w io.Writer, depth int) {
	n := rand.Intn(ch.high - ch.low)
	w.Write([]byte{byte(n + ch.low)})
}

type epsilon struct{}

func (e epsilon) generate(w io.Writer, depth int) {}
