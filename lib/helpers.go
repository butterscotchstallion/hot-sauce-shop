package lib

import (
	"fmt"
	"strings"
)

func debugQuery(query string, params []string) string {
	q := strings.ReplaceAll(query, "\n", " ")
	q = strings.ReplaceAll(q, "\t", " ")
	counter := 1
	for _, param := range params {
		q = strings.ReplaceAll(q, fmt.Sprintf("$%d", counter), fmt.Sprintf("'%s'", param))
		counter++
	}
	return q
}
