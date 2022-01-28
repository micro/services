package handler

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestLexer(t *testing.T) {

	tcs := []struct {
		name   string
		input  string
		tokens []item
		err    error
	}{
		{
			name:  "basic",
			input: `foo == "bar"`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `bar`,
				},
			},
		},
		{
			name:  "basic",
			input: `first_name == 'Dom'`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "first_name",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `Dom`,
				},
			},
		},
		{
			name:  "basic compressed",
			input: `foo=="bar"`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `bar`,
				},
			},
		},
		{
			name:  "basic bool",
			input: `foo == true`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemBoolean,
					val: `true`,
				},
			},
		},
		{
			name:  "basic bool false",
			input: `foo == false`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemBoolean,
					val: `false`,
				},
			},
		},
		{
			name:  "basic with spaces",
			input: `foo == "hello there"`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `hello there`,
				},
			},
		},
		{
			name:  "basic number",
			input: `foo == 123987`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemNumber,
					val: `123987`,
				},
			},
		},
		{
			name:  "basic gt number",
			input: `foo >= 123987`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: ">=",
				},
				{
					typ: itemNumber,
					val: `123987`,
				},
			},
		},
		{
			name:  "basic lt number",
			input: `foo <= 123987`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "<=",
				},
				{
					typ: itemNumber,
					val: `123987`,
				},
			},
		},
		{
			name:  "AND bool",
			input: `foo == 'bar' AND baz == 'hello'`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `bar`,
				},
				{
					typ: itemBooleanOp,
					val: "AND",
				},
				{
					typ: itemIdentifier,
					val: "baz",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `hello`,
				},
			},
		},
		{
			name:  "and bool",
			input: `foo == 'bar' and baz == 'hello'`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `bar`,
				},
				{
					typ: itemBooleanOp,
					val: "and",
				},
				{
					typ: itemIdentifier,
					val: "baz",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `hello`,
				},
			},
		},
		{
			name:  "OR bool",
			input: `foo == 'bar' OR baz == 'hello'`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `bar`,
				},
				{
					typ: itemBooleanOp,
					val: "OR",
				},
				{
					typ: itemIdentifier,
					val: "baz",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `hello`,
				},
			},
		},
		{
			name:  "or bool",
			input: `foo == 'bar' or baz == 'hello'`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `bar`,
				},
				{
					typ: itemBooleanOp,
					val: "or",
				},
				{
					typ: itemIdentifier,
					val: "baz",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `hello`,
				},
			},
		},
		{
			name:  "bad val",
			input: `foo == bar`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemError,
				},
			},
			err: fmt.Errorf("blah"),
		},
		{
			name:  "gibberish",
			input: `123onddlqkjn oajsldkj`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "123onddlqkjn",
				},
				{
					typ: itemError,
				},
			},
			err: fmt.Errorf("blah"),
		},
		{
			name:  "gibberish",
			input: `123onddlqkjn`,
			tokens: []item{
				{
					typ: itemError,
				},
			},
			err: fmt.Errorf("blah"),
		},
		{
			name:  "brackets",
			input: `foo == 'bar' and (baz == 'hello' or customer.name == 'john doe')`,
			tokens: []item{
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `bar`,
				},
				{
					typ: itemBooleanOp,
					val: "and",
				},
				{
					typ: itemLeftParen,
					val: "(",
				},
				{
					typ: itemIdentifier,
					val: "baz",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `hello`,
				},
				{
					typ: itemBooleanOp,
					val: "or",
				},
				{
					typ: itemIdentifier,
					val: "customer.name",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: "john doe",
				},
				{
					typ: itemRightParen,
					val: ")",
				},
			},
		},
		{
			name:  "brackets",
			input: `(foo == 'bar' and baz == 'hello') or customer.name == 'john doe'`,
			tokens: []item{
				{
					typ: itemLeftParen,
					val: "(",
				},
				{
					typ: itemIdentifier,
					val: "foo",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `bar`,
				},
				{
					typ: itemBooleanOp,
					val: "and",
				},
				{
					typ: itemIdentifier,
					val: "baz",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: `hello`,
				},
				{
					typ: itemRightParen,
					val: ")",
				},
				{
					typ: itemBooleanOp,
					val: "or",
				},
				{
					typ: itemIdentifier,
					val: "customer.name",
				},
				{
					typ: itemOperator,
					val: "==",
				},
				{
					typ: itemString,
					val: "john doe",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			_, ch := lex(tc.name, tc.input)
			erred := false
			for _, tok := range tc.tokens {
				it := <-ch
				t.Logf("Got %v", it)
				if it.typ == itemError {
					g.Expect(tok.typ).To(Equal(itemError))
					erred = true
				} else {
					g.Expect(it).To(Equal(tok))
				}

			}
			if tc.err != nil {
				g.Expect(erred).To(BeTrue())
			}
		})
	}

}
