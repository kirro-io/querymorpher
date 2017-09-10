package querymorpher

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	unquotedValueRegexp *regexp.Regexp
	operatorRegexp      *regexp.Regexp
)

const (
	ORDER_BY = "order_by"
	LIMIT    = "limit"
)

func init() {
	unquotedValueRegexp = regexp.MustCompile("(^[+-]?[0-9]+[.0-9]*$)|(^true|false$)")
	operatorRegexp = regexp.MustCompile("__([n]?eq|[gl]+t[e]?)$")
}

// queryScript is a builder for sql like query.
type queryScript struct {
	where   []string
	orderBy string
	limit   string
}

// set parses operator from key add quotes string value if needed. Parsed
// key and value are appended into WHERE clausule.
func (q *queryScript) set(key, value string) error {
	var fld, op, expr string

	op = operatorRegexp.FindString(key)
	op = getOperator(op)

	fld = operatorRegexp.ReplaceAllString(key, "")

	if fld == "" || value == "" {
		return fmt.Errorf("could not set '%s %s %s'", fld, op, value)
	}

	if strings.HasPrefix(value, "'") || strings.HasPrefix(value, `"`) || unquotedValueRegexp.MatchString(value) {
		expr = fmt.Sprintf("%s %s %s", fld, op, value)
	} else {
		expr = fmt.Sprintf("%s %s '%s'", fld, op, value)
	}

	q.where = append(q.where, expr)
	return nil
}

// setOrderBy sets `orderBy` attribute to given value.
func (q *queryScript) setOrderBy(value string) {
	if strings.HasPrefix(value, "-") {
		value = strings.Replace(value, "-", "", 1)
		q.orderBy = fmt.Sprintf("%s DESC", value)
	} else {
		q.orderBy = value
	}
}

// setLimit sets `limit` attribute to given value.
func (q *queryScript) setLimit(value string) {
	q.limit = value
}

// repr returns string representation of query.
func (q *queryScript) repr() (query string) {
	query = strings.Join(q.where, " AND ")

	if q.orderBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", q.orderBy)
	}

	if q.limit != "" {
		query += fmt.Sprintf(" LIMIT %s", q.limit)
	}

	return
}

// getOperator identifies and returns query operator.
func getOperator(op string) string {
	op = strings.Replace(op, "__", "", -1)
	switch op {
	case "eq":
		return "="
	case "neq":
		return "!="
	case "lt":
		return "<"
	case "lte":
		return "<="
	case "gt":
		return ">"
	case "gte":
		return ">="
	}
	return "="
}

// QueryFromRequest transform request's url query values to sql like query.
func QueryFromRequest(r *http.Request) (string, error) {
	query := &queryScript{}
	for key, val := range r.URL.Query() {
		if len(val) > 1 || len(val) == 0 {
			return "", fmt.Errorf("key '%s' cannot have none or multiple values", key)
		}
		value := val[0]
		switch key {
		case ORDER_BY:
			query.setOrderBy(value)
		case LIMIT:
			query.setLimit(value)
		default:
			err := query.set(key, value)
			if err != nil {
				return "", err
			}
		}
	}
	return query.repr(), nil
}
