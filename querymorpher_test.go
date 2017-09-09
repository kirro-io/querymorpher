package querymorpher

import (
	"fmt"
	"net/url"
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
		{"test gt operator", "number__gt=42", "number > 42", false},
		{"test lt operator", "number__lt=42", "number < 42", false},
		{"test gte operator", "number__gte=42", "number >= 42", false},
		{"test lte operator", "number__lte=42", "number <= 42", false},
		{"test non-existing operator", "t__t=test", "t__t = 'test'", false},
		{"test data", "date=2017-09-09", "date = '2017-09-09'", false},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			q, _ := url.ParseQuery(tc.query)
			res, err := Transform(q)
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

func ExampleTransform() {
	q, _ := url.ParseQuery("age__gte=18&name=John")
	res, _ := Transform(q)
	fmt.Println(res)
	// Output: age >= 18 AND name = 'John'
}
