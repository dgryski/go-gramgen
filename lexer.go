package main

import (
	"log"
)

var _scanner_actions = []int8{0, 1, 0, 1, 1, 1, 2, 1, 3, 1, 4, 1, 5, 1, 6, 1, 7, 1, 8, 0}
var _scanner_key_offsets = []int8{0, 0, 1, 2, 3, 16, 0}
var _scanner_trans_keys = []byte{34, 46, 61, 9, 10, 32, 34, 46, 58, 59, 95, 124, 65, 90, 97, 122, 95, 48, 57, 65, 90, 97, 122, 0}
var _scanner_single_lengths = []int8{0, 1, 1, 1, 9, 1, 0}
var _scanner_range_lengths = []int8{0, 0, 0, 0, 2, 3, 0}
var _scanner_index_offsets = []int8{0, 0, 2, 4, 6, 18, 0}
var _scanner_cond_targs = []int8{4, 1, 4, 0, 4, 0, 4, 4, 4, 1, 2, 3, 4, 5, 4, 5, 5, 0, 5, 5, 5, 5, 4, 0, 1, 2, 3, 4, 4, 0}
var _scanner_cond_actions = []int8{11, 0, 9, 0, 7, 0, 13, 15, 13, 0, 0, 0, 5, 0, 5, 0, 0, 0, 0, 0, 0, 0, 17, 0, 0, 0, 0, 0, 17, 0}
var _scanner_to_state_actions = []int8{0, 0, 0, 0, 1, 0, 0}
var _scanner_from_state_actions = []int8{0, 0, 0, 0, 3, 0, 0}
var _scanner_eof_trans = []int8{24, 25, 26, 27, 28, 29, 0}
var scanner_start int = 4
var _ = scanner_start
var scanner_first_final int = 4
var _ = scanner_first_final
var scanner_error int = 0
var _ = scanner_error
var scanner_en_main int = 4
var _ = scanner_en_main

type tok struct {
	t      int
	yy     yySymType
	lineno int
}

func lex(data []byte) *fuzzLexer {

	cs, p, pe, eof := 0, 0, len(data), len(data)
	ts, te, act := 0, 0, 0

	_, _, _ = ts, te, act

	lineno := 1

	var tokens []tok

	add := func(t int) {
		tokens = append(tokens, tok{t: t, lineno: lineno})
	}

	addstr := func(t int, s string) {
		tokens = append(tokens, tok{t: t, yy: yySymType{s: s}, lineno: lineno})
	}

	{
		cs = int(scanner_start)
		ts = 0
		te = 0
	}
	{
		var _klen int
		var _trans uint = 0
		var _keys int
		var _acts int
		var _nacts uint
	_resume:
		{
		}
		if p == pe && p != eof {
			goto _out
		}
		_acts = int(_scanner_from_state_actions[cs])

		_nacts = uint(_scanner_actions[_acts])
		_acts += 1
		for _nacts > 0 {
			switch _scanner_actions[_acts] {
			case 1:
				{
					{
						ts = p
					}
				}
			}
			_nacts -= 1
			_acts += 1
		}
		if p == eof {
			if _scanner_eof_trans[cs] > 0 {
				_trans = uint(_scanner_eof_trans[cs]) - 1
			}
		} else {
			_keys = int(_scanner_key_offsets[cs])

			_trans = uint(_scanner_index_offsets[cs])
			_klen = int(_scanner_single_lengths[cs])
			if _klen > 0 {
				var _lower int = _keys
				var _upper int = _keys + _klen - 1
				var _mid int
				for {
					if _upper < _lower {
						_keys += _klen
						_trans += uint(_klen)
						break
					}
					_mid = _lower + ((_upper - _lower) >> 1)
					if (data[p]) < _scanner_trans_keys[_mid] {
						_upper = _mid - 1
					} else if (data[p]) > _scanner_trans_keys[_mid] {
						_lower = _mid + 1
					} else {
						_trans += uint((_mid - _keys))
						goto _match
					}
				}
			}
			_klen = int(_scanner_range_lengths[cs])
			if _klen > 0 {
				var _lower int = _keys
				var _upper int = _keys + (_klen << 1) - 2
				var _mid int
				for {
					if _upper < _lower {
						_trans += uint(_klen)
						break
					}
					_mid = _lower + (((_upper - _lower) >> 1) & ^1)
					if (data[p]) < _scanner_trans_keys[_mid] {
						_upper = _mid - 2
					} else if (data[p]) > _scanner_trans_keys[_mid+1] {
						_lower = _mid + 2
					} else {
						_trans += uint(((_mid - _keys) >> 1))
						break
					}
				}
			}
		_match:
			{
			}
		}
		cs = int(_scanner_cond_targs[_trans])
		if _scanner_cond_actions[_trans] != 0 {
			_acts = int(_scanner_cond_actions[_trans])

			_nacts = uint(_scanner_actions[_acts])
			_acts += 1
			for _nacts > 0 {
				switch _scanner_actions[_acts] {
				case 2:
					{
						{
							te = p + 1
							{
								add(int(data[ts]))
							}
						}
					}
				case 3:
					{
						{
							te = p + 1
							{
								add(tASSIGN)
							}
						}
					}
				case 4:
					{
						{
							te = p + 1
							{
								add(tDOTDOT)
							}
						}
					}
				case 5:
					{
						{
							te = p + 1
							{
								addstr(tQSTRING, string(data[ts+1:te-1]))
							}
						}
					}
				case 6:
					{
						{
							te = p + 1
							{
							}
						}
					}
				case 7:
					{
						{
							te = p + 1
							{
								lineno++
							}
						}
					}
				case 8:
					{
						{
							te = p
							p = p - 1
							{
								addstr(tID, string(data[ts:te]))
							}
						}
					}
				}
				_nacts -= 1
				_acts += 1
			}
		}
		if p == eof {
			if cs >= 4 {
				goto _out
			}
		} else {
			_acts = int(_scanner_to_state_actions[cs])

			_nacts = uint(_scanner_actions[_acts])
			_acts += 1
			for _nacts > 0 {
				switch _scanner_actions[_acts] {
				case 0:
					{
						{
							ts = 0
						}
					}
				}
				_nacts -= 1
				_acts += 1
			}
			if cs != 0 {
				p += 1
				goto _resume
			}
		}
	_out:
		{
		}
	}
	return &fuzzLexer{toks: tokens}
}

type fuzzLexer struct {
	t    tok
	toks []tok
}

func (f *fuzzLexer) Lex(lval *yySymType) int {
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
