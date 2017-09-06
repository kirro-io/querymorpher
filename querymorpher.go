package querymorpher

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var operatorMap = map[string]string{
	"gt":  ">",
	"gte": ">=",
	"lt":  "<",
	"lte": "<=",
	"eq":  "=",
	"neq": "!=",
}

// Transform takes url.Values and transforms them to sql like query.
func Transform(u url.Values) (string, error) {
	var res []string

	for key, value := range u {
		if len(value) != 1 || len(value[0]) == 0 {
			return "", fmt.Errorf("multiple or no values are not supported for '%s'", key)
		}
		if len(key) == 0 {
			return "", fmt.Errorf("key for value '%s' cannot be empty", value[0])
		}

		attr, op := parseQueryKey(key)
		val := parseQueryValue(value[0])

		res = append(res, fmt.Sprintf("%s %s %s", attr, op, val))
	}

	return strings.Join(res, " AND "), nil
}

// parseQueryKey tries to get operator from query key. If operator
// is nout found in operatorMap then given key is returned with default
// operator "=".
func parseQueryKey(key string) (string, string) {

	s := strings.Split(key, "__")
	operator := s[len(s)-1]

	op, ok := operatorMap[operator]
	if !ok {
		return key, "="
	}

	return strings.Join(s[:len(s)-1], "__"), op
}

// parseQueryValue makes sure that string values are quoted.
func parseQueryValue(val string) string {
	ch := []rune(val)[0]
	if ch == '\'' || ch == '"' {
		return val
	}

	if _, err := strconv.ParseBool(val); err == nil {
		return val
	}
	if _, err := strconv.ParseInt(val, 10, 64); err == nil {
		return val
	}
	if _, err := strconv.ParseFloat(val, 64); err == nil {
		return val
	}

	return fmt.Sprintf("'%s'", val)
}
