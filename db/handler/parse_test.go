package handler

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/crufter/lexer"
)

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

	tokens, err = lexer.Lex(`a == 12 and name != "nandos"`, expressions)
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
	Q string
	E []Query
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
		// test escaping quotes
		tCase{
			Q: `a == 12 and name != "He said ""yes""!"`,
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
