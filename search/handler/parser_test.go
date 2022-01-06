package handler

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestParsing(t *testing.T) {

	tcs := []struct {
		name   string
		input  string
		output string
		err    error
	}{
		{
			name:   "basic",
			input:  `foo == "bar"`,
			output: `{"query":{"bool":{"must":[{"match":{"foo":"bar"}}]}}}`,
		},
		{
			name:   "basic",
			input:  `first_name == 'Dom'`,
			output: `{"query":{"bool":{"must":[{"match":{"first_name":"Dom"}}]}}}`,
		},
		{
			name:   "basic bool",
			input:  `foo == true`,
			output: `{"query":{"bool":{"must":[{"match":{"foo":"true"}}]}}}`,
		},
		{
			name:   "basic bool false",
			input:  `foo == false`,
			output: `{"query":{"bool":{"must":[{"match":{"foo":"false"}}]}}}`,
		},
		{
			name:   "basic with spaces",
			input:  `foo == "hello there"`,
			output: `{"query":{"bool":{"must":[{"match":{"foo":"hello there"}}]}}}`,
		},
		{
			name:   "basic number",
			input:  `foo == 123987`,
			output: `{"query":{"bool":{"must":[{"match":{"foo":"123987"}}]}}}`,
		},
		{
			name:   "basic gt number",
			input:  `foo >= 123987`,
			output: `{"query":{"bool":{"must":[{"range":{"foo":{"gte":"123987"}}}]}}}`,
		},
		{
			name:   "basic lt number",
			input:  `foo <= 123987`,
			output: `{"query":{"bool":{"must":[{"range":{"foo":{"lte":"123987"}}}]}}}`,
		},
		{
			name:   "AND bool",
			input:  `foo == 'bar' AND baz == 'hello'`,
			output: `{"query":{"bool":{"must":[{"match":{"foo":"bar"}},{"match":{"baz":"hello"}}]}}}`,
		},
		{
			name:   "and bool",
			input:  `foo == 'bar' and baz == 'hello'`,
			output: `{"query":{"bool":{"must":[{"match":{"foo":"bar"}},{"match":{"baz":"hello"}}]}}}`,
		},
		{
			name:   "OR bool",
			input:  `foo == 'bar' OR baz == 'hello'`,
			output: `{"query":{"bool":{"should":[{"match":{"foo":"bar"}},{"match":{"baz":"hello"}}]}}}`,
		},
		{
			name:   "or bool",
			input:  `foo == 'bar' or baz == 'hello'`,
			output: `{"query":{"bool":{"should":[{"match":{"foo":"bar"}},{"match":{"baz":"hello"}}]}}}`,
		},
		{
			name:  "bad val",
			input: `foo == bar`,
			err:   fmt.Errorf("blah"),
		},
		{
			name:  "gibberish",
			input: `123onddlqkjn oajsldkj`,
			err:   fmt.Errorf("blah"),
		},
		{
			name:  "gibberish",
			input: `123onddlqkjn`,
			err:   fmt.Errorf("blah"),
		},
		{
			name:   "brackets",
			input:  `foo == 'bar' and (baz == 'hello' or customer.name == 'john doe')`,
			output: `{"query":{"bool":{"must":[{"match":{"foo":"bar"}},{"bool":{"should":[{"match":{"baz":"hello"}},{"match":{"customer.name":"john doe"}}]}}]}}}`,
		},
		{
			name:   "gte",
			input:  `foo >= 6`,
			output: `{"query":{"bool":{"must":[{"range":{"foo":{"gte":"6"}}}]}}}`,
		},
		{
			name:   "lte",
			input:  `foo <= 6`,
			output: `{"query":{"bool":{"must":[{"range":{"foo":{"lte":"6"}}}]}}}`,
		},
		{
			name:   "wildcard",
			input:  `foo == "ba*"`,
			output: `{"query":{"bool":{"must":[{"wildcard":{"foo":{"value":"ba*"}}}]}}}`,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			js, err := parseQueryString(tc.input)
			if tc.err != nil {
				g.Expect(err).To(Not(BeNil()))
			} else {
				b, _ := js.MarshalJSON()
				t.Logf("%+v", string(b))
				g.Expect(err).To(BeNil())
				g.Expect(string(b)).To(Equal(tc.output))
			}

		})
	}

}
