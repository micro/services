package handler

import (
	"regexp"
)

var (
	IDFormat = regexp.MustCompilePOSIX("[a-z0-9-]+")

	FunctionKey = "function/func/"
	OwnerKey = "function/owner/"
)
