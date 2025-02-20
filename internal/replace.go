package internal

import "regexp"

func Replace(content string, expression string, replacement string) string {
	re := regexp.MustCompile(expression)
	return re.ReplaceAllString(content, replacement)
}
