package handler

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/crufter/lexer"
)

func TestCorrectFieldName(t *testing.T) {
	f := correctFieldName("a.b.c", true)
	if f != "data -> 'a' -> 'b' ->> 'c'" {
		t.Fatal(f)
	}
}

func TestLexing(t *testing.T) {
	tokens, err := lexer.Lex("a == 12", expressions)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 3 {
		t.Fatal(tokens)
	}
	if tokens[0].Typ != itemFieldName || tokens[1].Typ != itemEquals || tokens[2].Typ != itemInt {
		t.Fatal(tokens)
	}

	tokens, err = lexer.Lex(`a == 12 and name != 'nandos'`, expressions)
	if tokens[0].Typ != itemFieldName ||
		tokens[1].Typ != itemEquals ||
		tokens[2].Typ != itemInt ||
		tokens[3].Typ != itemAnd ||
		tokens[4].Typ != itemFieldName ||
		tokens[5].Typ != itemNotEquals ||
		tokens[6].Typ != itemString {
		t.Fatal(tokens)
	}
}

type tCase struct {
	Q   string
	E   []Query
	Err error
}

func TestParsing(t *testing.T) {
	tCases := []tCase{
		tCase{
			Q: `a == 12 and name != "nandos"`,
			E: []Query{
				Query{
					Field: "a",
					Value: int64(12),
					Op:    itemEquals,
				},
				Query{
					Field: "name",
					Value: "nandos",
					Op:    itemNotEquals,
				},
			},
		},
		tCase{
			Q: `a.b.c == 12 and name != "nandos"`,
			E: []Query{
				Query{
					Field: "a.b.c",
					Value: int64(12),
					Op:    itemEquals,
				},
				Query{
					Field: "name",
					Value: "nandos",
					Op:    itemNotEquals,
				},
			},
		},
		tCase{
			Q: `a == 12 and name != "nan'dos"`,
			E: []Query{
				Query{
					Field: "a",
					Value: int64(12),
					Op:    itemEquals,
				},
				Query{
					Field: "name",
					Value: "nan'dos",
					Op:    itemNotEquals,
				},
			},
		},
		tCase{
			Q: `id == '795c1e56-d1f3-495d-b9cb-d84a56ffb39c'`,
			E: []Query{
				Query{
					Field: "id",
					Value: "795c1e56-d1f3-495d-b9cb-d84a56ffb39c",
					Op:    itemEquals,
				},
			},
		},
		tCase{
			Q: `a == 12 and name != 'nandos'`,
			E: []Query{
				Query{
					Field: "a",
					Value: int64(12),
					Op:    itemEquals,
				},
				Query{
					Field: "name",
					Value: "nandos",
					Op:    itemNotEquals,
				},
			},
		},
		tCase{
			Q: "a == 12 and name != `nandos`",
			E: []Query{
				Query{
					Field: "a",
					Value: int64(12),
					Op:    itemEquals,
				},
				Query{
					Field: "name",
					Value: "nandos",
					Op:    itemNotEquals,
				},
			},
		},
		// test escaping quotes
		tCase{
			Q: `a == 12 and name != 'He said ""yes""!'`,
			E: []Query{
				Query{
					Field: "a",
					Value: int64(12),
					Op:    itemEquals,
				},
				Query{
					Field: "name",
					Value: `He said "yes"!`,
					Op:    itemNotEquals,
				},
			},
		},
		tCase{
			Q: `a == false and b == true`,
			E: []Query{
				Query{
					Field: "a",
					Value: false,
					Op:    itemEquals,
				},
				Query{
					Field: "b",
					Value: true,
					Op:    itemEquals,
				},
			},
		},
		// a < 20
		tCase{
			Q: `a < 20`,
			E: []Query{
				Query{
					Field: "a",
					Value: int64(20),
					Op:    itemLessThan,
				},
			},
		},
		tCase{
			Q: `a <= 20`,
			E: []Query{
				Query{
					Field: "a",
					Value: int64(20),
					Op:    itemLessThanEquals,
				},
			},
		},
		tCase{
			Q: `a > 20`,
			E: []Query{
				Query{
					Field: "a",
					Value: int64(20),
					Op:    itemGreaterThan,
				},
			},
		},
		tCase{
			Q: `a >= 20`,
			E: []Query{
				Query{
					Field: "a",
					Value: int64(20),
					Op:    itemGreaterThanEquals,
				},
			},
		},
	}
	for _, tCase := range tCases {
		fmt.Println("Parsing", tCase.Q)
		qs, err := Parse(tCase.Q)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(qs, tCase.E) {
			t.Fatal("Expected", tCase.E, "got", qs)
		}
	}
}
