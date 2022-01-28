package handler

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestEmailValidation(t *testing.T) {
	tcs := []struct {
		name  string
		email string
		valid bool
	}{
		{
			name:  "Empty",
			email: "",
			valid: false,
		},
		{
			name:  "Normal email",
			email: "joe@example.com",
			valid: true,
		},
		{
			name:  "Email with dots",
			email: "joe.bloggs@example.com",
			valid: true,
		},
		{
			name:  "Email with plus",
			email: "joe+1@example.com",
			valid: true,
		},
		{
			name:  "Email with underscores",
			email: "joe_bloggs@example.com",
			valid: true,
		},
		{
			name:  "White space",
			email: "joe bloggs@example.com",
			valid: false,
		},
		{
			name:  "Trailing whitespace",
			email: "joe@example.com ",
			valid: false,
		},
		{
			name:  "Preceding whitespace",
			email: " joe@example.com",
			valid: false,
		},
		{
			name:  "Normal email",
			email: "joe@example.com",
			valid: true,
		},
	}
	for _, tc := range tcs {
		g := NewWithT(t)
		t.Run(tc.name, func(t *testing.T) {
			g.Expect(validEmail(tc.email)).To(Equal(tc.valid))
		})
	}
}
