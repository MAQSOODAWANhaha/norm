// builder/query_test.go
package builder

import (
	"testing"

	"norm/model"
)

func TestQueryBuilder_Validation(t *testing.T) {
	registry := model.NewRegistry()
	qb := NewQueryBuilder(registry)

	// Test case 1: Valid query
	t.Run("Valid Query", func(t *testing.T) {
		result, err := qb.Match("(n:Person)").Return("n").Build()
		if err != nil {
			t.Fatalf("Build failed unexpectedly: %v", err)
		}
		if !result.Valid {
			t.Errorf("Expected query to be valid, but it was invalid. Errors: %v", result.Errors)
		}
		if len(result.Errors) > 0 {
			t.Errorf("Expected no errors, but got %d", len(result.Errors))
		}
	})

	// Test case 2: Invalid query (bracket mismatch)
	t.Run("Invalid Query with Mismatched Brackets", func(t *testing.T) {
		// Create a new builder for a clean state
		qb2 := NewQueryBuilder(registry)
		result, err := qb2.Match("(n:Person").Return("n").Build()
		if err != nil {
			t.Fatalf("Build failed unexpectedly: %v", err)
		}
		if result.Valid {
			t.Error("Expected query to be invalid, but it was valid.")
		}
		if len(result.Errors) == 0 {
			t.Error("Expected errors, but got none.")
		} else {
			if result.Errors[0].Type != "bracket_mismatch" {
				t.Errorf("Expected error type 'bracket_mismatch', but got '%s'", result.Errors[0].Type)
			}
		}
	})
}
