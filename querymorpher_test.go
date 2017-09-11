package querymorpher

import (
	"net/http"
	"testing"
)

func TestTransform(t *testing.T) {
	tt := []struct {
		name       string
		query      string
		expected   string
		shouldFail bool
	}{
		{"test empty query", "", "", false},
		{"test bool", "fld=true", "fld = true", false},
		{"test float", "fld=4.2", "fld = 4.2", false},
		{"test int", "fld=4", "fld = 4", false},
		{"test string with quotes", "fld='test'", "fld = 'test'", false},
		{"test string without quotes", "fld=test", "fld = 'test'", false},
		{"test empty key", "=test", "", true},
		{"test empty value", "key=", "", true},
		{"test duplicate key", "key=1&key=1", "", true},
		{"test eq operator", "number__eq=42", "number = 42", false},
		{"test neq operator", "number__neq=42", "number != 42", false},
		{"test gt operator", "number__gt=42", "number > 42", false},
		{"test lt operator", "number__lt=42", "number < 42", false},
		{"test gte operator", "number__gte=42", "number >= 42", false},
		{"test lte operator", "number__lte=42", "number <= 42", false},
		{"test non-existing operator", "t__t=test", "t__t = 'test'", false},
		{"test date", "date=2017-09-09", "date = '2017-09-09'", false},
		{"test limit", "age=18&limit=1", "age = 18 LIMIT 1", false},
		{"test order by", "age__gt=18&order_by=age", "age > 18 ORDER BY age", false},
		{"test order by desc", "age__gt=18&order_by=-age", "age > 18 ORDER BY age DESC", false},
		{"test order by + limit", "age__gt=18&order_by=age&limit=2", "age > 18 ORDER BY age LIMIT 2", false},
		{"test order by desc + limit", "age__gt=18&order_by=-age&limit=2", "age > 18 ORDER BY age DESC LIMIT 2", false},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", "?"+tc.query, nil)
			res, err := QueryFromRequest(r)
			if !tc.shouldFail && err != nil {
				t.Error(err.Error())
			}
			if tc.shouldFail && err == nil {
				t.Errorf("should fail for '%s' but got '%s'", tc.query, res)
			}
			if res != tc.expected {
				t.Errorf("Expected: '%s', Got: '%s'", tc.expected, res)
			}
		})
	}
}
