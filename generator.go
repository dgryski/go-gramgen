package main

import (
	"io"
	"log"
	"math/rand"
	"strconv"
)

type generator interface {
	generate(w io.Writer)
}

type terminal string

func (t terminal) generate(w io.Writer) {
	w.Write([]byte(t))
}

type variable string

func (v variable) generate(w io.Writer) {
	g, ok := symtab[string(v)]
	if !ok {
		log.Fatalf("unknown variable %q", string(v))
	}
	g.generate(w)
}

type choice []generator

func (c choice) generate(w io.Writer) {
	n := rand.Intn(len(c))
	c[n].generate(w)
}

func (c *choice) add(g generator) {
	*c = append(*c, g)
}

type sequence []generator

func (s sequence) generate(w io.Writer) {
	for _, ss := range s {
		ss.generate(w)
	}
}

func (s *sequence) add(g generator) {
	*s = append(*s, g)
}

type intrange struct {
	low, high int
}

func (ir intrange) generate(w io.Writer) {
	n := rand.Intn(ir.high - ir.low)
	w.Write(strconv.AppendInt(nil, int64(ir.low+n), 10))
}

type chrange struct {
	low, high int
}

func (ch chrange) generate(w io.Writer) {
	n := rand.Intn(ch.high - ch.low)
	w.Write([]byte{byte(n + ch.low)})
}

type epsilon struct{}

func (e epsilon) generate(w io.Writer) {}
