package server

import (
	"fmt"
	models "querylang-chart/graph/model"
	"strings"
)

func DeserializeQuery(query string) models.DeserializedQuery {
	var dq = models.DeserializedQuery{}
	if strings.Count(query, ")") == 1 {
		dq.Query = translate(query)
	} else {
		parenthesesIndex := strings.Index(query, "(")
		queryKeyWord := keyWord(query, parenthesesIndex)

		if queryKeyWord == "NOT" {
			dq.Query = fmt.Sprintf("(%s %s)", queryKeyWord, DeserializeQuery(query[parenthesesIndex+1:len(query)-1]).Query)
		} else {
			level := countClosingParentheses(query)
			commaIndex := strings.Index(query, strings.Repeat(")", level-1)+",") + level - 1
			dq.Query = fmt.Sprintf("(%s %s %s)", DeserializeQuery(query[parenthesesIndex+1:commaIndex]).Query, queryKeyWord, DeserializeQuery(query[commaIndex+1:len(query)-1]).Query)
		}
	}
	return dq
}

func translate(s string) string {
	parenthesesIndex := strings.Index(s, "(")
	commaIndex := strings.Index(s, ",")
	keyWord := keyWord(s, parenthesesIndex)
	return s[parenthesesIndex+1:commaIndex] + keyWord + s[commaIndex+1:len(s)-1]
}

func keyWord(s string, i int) string {
	keyWord := s[:i]
	switch keyWord {
	case "EQUAL":
		return "="
	case "GREATER_THAN":
		return ">"
	case "LESS_THAN":
		return "<"
	default:
		return keyWord
	}
}

func countClosingParentheses(s string) int {
	counter := 0
	for i := len(s) - 1; i >= 0; i-- {
		if string(s[i]) == ")" {
			counter++
		} else {
			break
		}
	}
	return counter
}
