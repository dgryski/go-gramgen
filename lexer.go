
//line lexer.rl:1
package main

import (
	"fmt"
)


//line lexer.rl:8

//line lexer.go:13
var _scanner_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 1, 3, 
	1, 4, 1, 5, 1, 6, 1, 7, 
	1, 8, 
}

var _scanner_key_offsets []byte = []byte{
	0, 0, 1, 2, 3, 16, 
}

var _scanner_trans_keys []byte = []byte{
	34, 46, 61, 9, 10, 32, 34, 46, 
	58, 59, 95, 124, 65, 90, 97, 122, 
	95, 48, 57, 65, 90, 97, 122, 
}

var _scanner_single_lengths []byte = []byte{
	0, 1, 1, 1, 9, 1, 
}

var _scanner_range_lengths []byte = []byte{
	0, 0, 0, 0, 2, 3, 
}

var _scanner_index_offsets []byte = []byte{
	0, 0, 2, 4, 6, 18, 
}

var _scanner_trans_targs []byte = []byte{
	4, 1, 4, 0, 4, 0, 4, 4, 
	4, 1, 2, 3, 4, 5, 4, 5, 
	5, 0, 5, 5, 5, 5, 4, 4, 
	
}

var _scanner_trans_actions []byte = []byte{
	11, 0, 9, 0, 7, 0, 13, 15, 
	13, 0, 0, 0, 5, 0, 5, 0, 
	0, 0, 0, 0, 0, 0, 17, 17, 
	
}

var _scanner_to_state_actions []byte = []byte{
	0, 0, 0, 0, 1, 0, 
}

var _scanner_from_state_actions []byte = []byte{
	0, 0, 0, 0, 3, 0, 
}

var _scanner_eof_trans []byte = []byte{
	0, 0, 0, 0, 0, 24, 
}

const scanner_start int = 4
const scanner_first_final int = 4
const scanner_error int = 0

const scanner_en_main int = 4


//line lexer.rl:9

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

	
//line lexer.go:102
	{
	cs = scanner_start
	ts = 0
	te = 0
	act = 0
	}

//line lexer.go:110
	{
	var _klen int
	var _trans int
	var _acts int
	var _nacts uint
	var _keys int
	if p == pe {
		goto _test_eof
	}
	if cs == 0 {
		goto _out
	}
_resume:
	_acts = int(_scanner_from_state_actions[cs])
	_nacts = uint(_scanner_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		 _acts++
		switch _scanner_actions[_acts - 1] {
		case 1:
//line NONE:1
ts = p

//line lexer.go:133
		}
	}

	_keys = int(_scanner_key_offsets[cs])
	_trans = int(_scanner_index_offsets[cs])

	_klen = int(_scanner_single_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + _klen - 1)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + ((_upper - _lower) >> 1)
			switch {
			case data[p] < _scanner_trans_keys[_mid]:
				_upper = _mid - 1
			case data[p] > _scanner_trans_keys[_mid]:
				_lower = _mid + 1
			default:
				_trans += int(_mid - int(_keys))
				goto _match
			}
		}
		_keys += _klen
		_trans += _klen
	}

	_klen = int(_scanner_range_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + (_klen << 1) - 2)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + (((_upper - _lower) >> 1) & ^1)
			switch {
			case data[p] < _scanner_trans_keys[_mid]:
				_upper = _mid - 2
			case data[p] > _scanner_trans_keys[_mid + 1]:
				_lower = _mid + 2
			default:
				_trans += int((_mid - int(_keys)) >> 1)
				goto _match
			}
		}
		_trans += _klen
	}

_match:
_eof_trans:
	cs = int(_scanner_trans_targs[_trans])

	if _scanner_trans_actions[_trans] == 0 {
		goto _again
	}

	_acts = int(_scanner_trans_actions[_trans])
	_nacts = uint(_scanner_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		_acts++
		switch _scanner_actions[_acts-1] {
		case 2:
//line lexer.rl:37
te = p+1
{ add(int(data[ts])) }
		case 3:
//line lexer.rl:38
te = p+1
{ add(tASSIGN) }
		case 4:
//line lexer.rl:39
te = p+1
{ add(tDOTDOT) }
		case 5:
//line lexer.rl:40
te = p+1
{ addstr(tQSTRING, string(data[ts+1:te-1])) }
		case 6:
//line lexer.rl:41
te = p+1
{ }
		case 7:
//line lexer.rl:42
te = p+1
{ lineno++ }
		case 8:
//line lexer.rl:43
te = p
p--
{
			addstr(tID, string(data[ts:te]))
		}
//line lexer.go:233
		}
	}

_again:
	_acts = int(_scanner_to_state_actions[cs])
	_nacts = uint(_scanner_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		_acts++
		switch _scanner_actions[_acts-1] {
		case 0:
//line NONE:1
ts = 0

//line lexer.go:247
		}
	}

	if cs == 0 {
		goto _out
	}
	p++
	if p != pe {
		goto _resume
	}
	_test_eof: {}
	if p == eof {
		if _scanner_eof_trans[cs] > 0 {
			_trans = int(_scanner_eof_trans[cs] - 1)
			goto _eof_trans
		}
	}

	_out: {}
	}

//line lexer.rl:50


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