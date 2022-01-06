package handler

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestParsing(t *testing.T) {

	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "basic",
			input: `foo == "bar"`,
		},
		{
			name:  "basic",
			input: `first_name == 'Dom'`,
		},
		{
			name:  "basic bool",
			input: `foo == true`,
		},
		{
			name:  "basic bool false",
			input: `foo == false`,
		},
		{
			name:  "basic with spaces",
			input: `foo == "hello there"`,
		},
		{
			name:  "basic number",
			input: `foo == 123987`,
		},
		{
			name:  "basic gt number",
			input: `foo >= 123987`,
		},
		{
			name:  "basic lt number",
			input: `foo <= 123987`,
		},
		{
			name:  "AND bool",
			input: `foo == 'bar' AND baz == 'hello'`,
		},
		{
			name:  "and bool",
			input: `foo == 'bar' and baz == 'hello'`,
		},
		{
			name:  "OR bool",
			input: `foo == 'bar' OR baz == 'hello'`,
		},
		{
			name:  "or bool",
			input: `foo == 'bar' or baz == 'hello'`,
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
			name:  "brackets",
			input: `foo == 'bar' and (baz == 'hello' or customer.name == 'john doe')`,
		},
		{
			name:  "gte",
			input: `foo >= 6`,
		},
		{
			name:  "lte",
			input: `foo <= 6`,
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
			}

		})
	}

}
