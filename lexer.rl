package main

import (
	"log"
)

%% machine scanner;
%% write data;

type tok struct {
    t int
    yy yySymType
    lineno int
}

func lex(data []byte) *fuzzLexer {

	cs, p, pe, eof := 0, 0, len(data), len(data)
	ts, te, act := 0, 0, 0

	_, _, _ = ts, te, act

	lineno := 1

	var tokens []tok

	add := func(t int) {
	    tokens = append(tokens, tok{t:t, lineno:lineno})
	}

	addstr := func(t int, s string) {
	    tokens = append(tokens, tok{t:t, yy:yySymType{s:s}, lineno:lineno})
	}

	%%{

	    main := |*
		[;|] => { add(int(data[ts])) };
		':=' => { add(tASSIGN) };
		'\.\.' => { add(tDOTDOT) };
		'"' [^"]* '"'  => { addstr(tQSTRING, string(data[ts+1:te-1])) };
		[ \t] => { };
		'\n' => { lineno++ };
		[A-Za-z_][A-Za-z0-9_]* {
			addstr(tID, string(data[ts:te]))
		};
	    *|;

	    write init;
	    write exec;
	}%%

	return &fuzzLexer{toks:tokens}
}

type fuzzLexer struct {
	t tok
	toks []tok
}

func (f *fuzzLexer) Lex(lval *yySymType) int  {
    if len(f.toks) == 0 {
	    return 0
    }
    f.t = f.toks[0]
    f.toks = f.toks[1:]
    *lval = f.t.yy
    return f.t.t
}

func (f *fuzzLexer) Error(s string) {
    log.Fatalf("syntax error at line %d: %v\n", f.t.lineno, s)
}
