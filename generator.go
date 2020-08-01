package main

import (
	"io"
	"strconv"
)

var xrand = xorm(1)

type generator interface {
	generate(w io.StringWriter, depth int)
}

type terminal string

func (t terminal) generate(w io.StringWriter, depth int) {
	w.WriteString(string(t))
}

type variable struct {
	v string
}

func (v variable) generate(w io.StringWriter, depth int) {
	if depth <= 0 {
		cheapestOption[v.v].generate(w, 0)
		return
	}
	g := symtab[v.v]
	g.generate(w, depth-1)
}

type choice struct {
	c     []generator
	cheap generator
}

func (c *choice) generate(w io.StringWriter, depth int) {
	if depth <= 0 {
		c.cheap.generate(w, depth-1)
		return
	}
	n := xrand.Intn(len(c.c))
	c.c[n].generate(w, depth-1)
}

func (c *choice) add(g generator) {
	c.c = append(c.c, g)
}

type sequence struct {
	s []generator
}

func (s *sequence) generate(w io.StringWriter, depth int) {
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

func (ir intrange) generate(w io.StringWriter, depth int) {
	n := xrand.Intn(ir.high - ir.low)
	w.WriteString(strconv.FormatInt(int64(ir.low+int(n)), 10))
}

type chrange struct {
	low, high int
}

func (ch chrange) generate(w io.StringWriter, depth int) {
	n := xrand.Intn(ch.high - ch.low)
	w.WriteString(string(rune(int(n) + ch.low)))
}

type epsilon struct{}

func (e epsilon) generate(w io.StringWriter, depth int) {}

type xorm uint64

func (r *xorm) Next() uint64 {
	x := *r
	x ^= x >> 12 // a
	x ^= x << 25 // b
	x ^= x >> 27 // c
	*r = x * 2685821657736338717
	return uint64(*r)
}

func (r *xorm) Intn(n int) uint64 {
	bound := uint64(n)
	threshold := -bound % bound
	for {
		n := r.Next()
		if n >= threshold {
			return n % bound
		}
	}
}
