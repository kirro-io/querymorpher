package querymorpher

import (
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
		{"test multiple strings", "fld=test&fld2='test'", "fld = 'test' AND fld2 = 'test'", false},
		{"test empty key", "=test", "", true},
		{"test empty value", "key=", "", true},
		{"test gt operator", "number__gt=42", "number > 42", false},
		{"test lt operator", "number__lt=42", "number < 42", false},
		{"test gte operator", "number__gte=42", "number >= 42", false},
		{"test lte operator", "number__lte=42", "number <= 42", false},
		{"test non-existing operator", "t__t=test", "t__t = 'test'", false},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			q, _ := url.ParseQuery(tc.query)
			res, err := Transform(q)
			if !tc.shouldFail && err != nil {
				t.Error(err.Error())
			}
			if tc.shouldFail && err == nil {
				t.Errorf("should fail for '%s' but got ''", tc.query, res)
			}
			if res != tc.expected {
				t.Errorf("Expected: '%s', Got: '%s'", res, tc.expected)
			}
		})
	}
}
