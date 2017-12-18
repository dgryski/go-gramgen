package main

import (
	"fmt"
)

%% machine scanner;
%% write data;

type tok struct {
    t int
    yy yySymType
}

func lex(data []byte) []tok {

	cs, p, pe, eof := 0, 0, len(data), len(data)
	ts, te, act := 0, 0, 0

	_, _, _ = ts, te, act

	lineno := 1

	var tokens []tok

	add := func(t int) {
	    tokens = append(tokens, tok{t:t})
	}

	addstr := func(t int, s string) {
	    tokens = append(tokens, tok{t:t, yy:yySymType{s:s}})
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

	return tokens
}

type fuzzLexer []tok

func (f *fuzzLexer) Lex(lval *yySymType) int  {
    if len(*f) == 0 {
	    return 0
    }
    t := (*f)[0]
    *f = (*f)[1:]
    *lval = t.yy
    return t.t
}

func (f *fuzzLexer) Error(s string) {
    fmt.Println("syntax error:", s)
}