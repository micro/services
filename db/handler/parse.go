package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/crufter/lexer"
)

var quoteEscape = fmt.Sprint(0x10FFFF)
var singleQuoteEscape = fmt.Sprint(0x10FFFE)
var backtickEscape = fmt.Sprint(0x10FFFD)

const (
	itemIgnore = iota

	// nouns
	itemAnd
	itemInt
	itemFieldName
	itemString
	itemBoolTrue
	itemBoolFalse

	// ops
	itemEquals
	itemNotEquals
	itemLessThan
	itemGreaterThan
	itemLessThanEquals
	itemGreaterThanEquals
)

var opToString = map[int]string{
	itemEquals:            "==",
	itemNotEquals:         "!=",
	itemLessThan:          "<",
	itemGreaterThan:       ">",
	itemLessThanEquals:    "<=",
	itemGreaterThanEquals: ">=",
}

var expressions = []lexer.TokenExpr{
	{`[ ]+`, itemIgnore}, // Whitespace
	{`==`, itemEquals},
	{`!=`, itemNotEquals},
	{`false`, itemBoolFalse},
	{`true`, itemBoolTrue},
	{`and`, itemAnd},
	{`<=`, itemLessThanEquals},
	{`>=`, itemGreaterThanEquals},
	{`<`, itemLessThan},
	{`>`, itemGreaterThan},
	{`[0-9]+`, itemInt},
	{`"(?:[^"\\]|\\.)*"`, itemString},
	{"`" + `(?:[^"\\]|\\.)*` + "`", itemString},
	{`'(?:[^"\\]|\\.)*'`, itemString},
	{`[\<\>\!\=\+\-\|\&\*\/A-Za-z][A-Za-z0-9_\.]*`, itemFieldName},
}

type Query struct {
	Field string
	Op    int
	Value interface{}
}

func Parse(q string) ([]Query, error) {
	if strings.Contains(q, quoteEscape) {
		return nil, errors.New("query contains illegal max rune")
	}
	q = strings.Replace(q, `""`, quoteEscape, -1)
	q = strings.Replace(q, "``", singleQuoteEscape, -1)
	q = strings.Replace(q, "''", backtickEscape, -1)

	tokens, err := lexer.Lex(q, expressions)
	if err != nil {
		return nil, err
	}
	queries := []Query{}
	current := Query{}
	for i, token := range tokens {
		// and tokens should trigger a query
		// save and reset
		if token.Typ == itemAnd {
			queries = append(queries, current)
			current = Query{}
			continue
		}

		// is an op
		if token.Typ >= itemEquals {
			current.Op = token.Typ
			continue
		}

		// is a value
		switch token.Typ {
		case itemFieldName:
			current.Field = token.Text
		case itemString:
			switch current.Op {
			case itemEquals, itemNotEquals:
			default:
				return nil, fmt.Errorf("operator '%v' can't be used with strings", opToString[current.Op])
			}

			if len(token.Text) < 2 {
				return nil, fmt.Errorf("string literal too short: '%v'", token.Text)
			}
			to := token.Text[1 : len(token.Text)-1]
			to = strings.Replace(to, quoteEscape, `"`, -1)
			to = strings.Replace(to, singleQuoteEscape, `'`, -1)
			to = strings.Replace(to, backtickEscape, "`", -1)
			current.Value = to
		case itemBoolTrue:
			switch current.Op {
			case itemEquals, itemNotEquals:
			default:
				return nil, fmt.Errorf("operator '%v' can't be used with bools", opToString[current.Op])
			}
			current.Value = true
		case itemBoolFalse:
			switch current.Op {
			case itemEquals, itemNotEquals:
			default:
				return nil, fmt.Errorf("operator '%v' can't be used with bools", opToString[current.Op])
			}
			current.Value = false
		case itemInt:
			num, err := strconv.ParseInt(token.Text, 10, 64)
			if err != nil {
				return nil, err
			}
			current.Value = num
		}

		// if we are at last position, save last query
		if i == len(tokens)-1 {
			queries = append(queries, current)
		}
	}
	return queries, nil
}
