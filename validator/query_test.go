// validator/query_test.go
package validator

import (
	"testing"
)

func TestCypherQueryValidator_Validate(t *testing.T) {
	validator := NewQueryValidator(true)

	testCases := []struct {
		name          string
		query         string
		expectError   bool
		errorType     string
	}{
		{
			name:        "Valid Query",
			query:       "MATCH (n:Person) RETURN n",
			expectError: false,
		},
		{
			name:        "Empty Query",
			query:       "",
			expectError: true,
			errorType:   "empty_query",
		},
		{
			name:        "Mismatched Parentheses",
			query:       "MATCH (n:Person RETURN n",
			expectError: true,
			errorType:   "bracket_mismatch",
		},
		{
			name:        "Mismatched Square Brackets",
			query:       "MATCH (n:Person {name: 'test']) RETURN n",
			expectError: true,
			errorType:   "bracket_mismatch",
		},
		{
			name:        "No Valid Clause",
			query:       "(n:Person)",
			expectError: true,
			errorType:   "no_valid_clause",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errors := validator.Validate(tc.query)
			if tc.expectError {
				if len(errors) == 0 {
					t.Errorf("Expected validation errors but got none")
				}
				found := false
				for _, err := range errors {
					if err.Type == tc.errorType {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error type %s but got different errors", tc.errorType)
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("Expected no validation errors but got %v", errors)
				}
			}
		})
	}
}
