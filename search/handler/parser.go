package handler

import (
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"
)

func parseQueryString(query string) (*simplejson.Json, error) {
	_, items := lex("foo", query)
	res, err := parseQueryStringRec(items)
	if err != nil {
		return nil, err
	}
	js := simplejson.New()
	js.Set("query", res)
	return js, nil
}

const (
	matchTypeRange    = "range"
	matchTypeMatch    = "match"
	matchTypeWildcard = "wildcard"
)

func parseQueryStringRec(items chan item) (*simplejson.Json, error) {
	retTerm := simplejson.New()
	currFieldName := ""
	currBool := ""
	currMatchType := matchTypeMatch
	currPathAddition := ""
	terms := []*simplejson.Json{}
itemLoop:
	for it := range items {
		if it.typ == itemError {
			return nil, fmt.Errorf(it.val)
		}

		switch it.typ {
		case itemIdentifier:
			currFieldName = it.val
		case itemString, itemBoolean, itemNumber:
			currTerm := simplejson.New()
			if strings.ContainsRune(it.val, '*') {
				currMatchType = matchTypeWildcard
				currPathAddition = "value"
			}
			path := []string{currMatchType, currFieldName}
			if len(currPathAddition) > 0 {
				path = append(path, currPathAddition)
			}
			currTerm.SetPath(path, it.val)
			terms = append(terms, currTerm)

			// reset
			currFieldName = ""
			currMatchType = matchTypeMatch

		case itemBooleanOp:
			thisBool := "must"
			if strings.ToLower(it.val) == "or" {
				thisBool = "should"
			}
			if len(currBool) == 0 {
				currBool = thisBool
			} else {
				if thisBool != currBool {
					// undefined behaviour order of precedence needs to be explicitly defined using parentheses
					return nil, fmt.Errorf("error parsing query. Cannot mix boolean operators without explicitly defining precedence with parentheses")
				}
			}
		case itemLeftParen:
			currTerm, err := parseQueryStringRec(items)
			if err != nil {
				return nil, err
			}
			terms = append(terms, currTerm)
		case itemRightParen:
			break itemLoop
		case itemOperator:
			switch it.val {
			case "==":
				currMatchType = matchTypeMatch
				currPathAddition = ""
			case ">=":
				currMatchType = matchTypeRange
				currPathAddition = "gte"
			case "<=":
				currMatchType = matchTypeRange
				currPathAddition = "lte"
			}

		}
	}
	if len(currBool) == 0 {
		currBool = "must"
	}
	retTerm.SetPath([]string{"bool", currBool}, terms)
	return retTerm, nil
}
