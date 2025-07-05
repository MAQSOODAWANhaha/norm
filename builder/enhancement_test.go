// builder/enhancement_test.go
package builder

import (
	"testing"

	"norm/types"
)

func TestVariableLengthRelationship(t *testing.T) {
	t.Run("Variable Length Relationship with Min and Max", func(t *testing.T) {
		qb := NewQueryBuilder()
		min := 2
		max := 4
		pattern := types.Pattern{
			StartNode: types.NodePattern{Variable: "a"},
			Relationship: types.RelationshipPattern{
				Type:      "KNOWS",
				MinLength: &min,
				MaxLength: &max,
				Direction: types.DirectionOutgoing,
			},
			EndNode: types.NodePattern{Variable: "b"},
		}

		result, err := qb.MatchPattern(pattern).Return("a, b").Build()
		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		expectedQuery := "MATCH (a)-[:KNOWS*2..4]->(b)\nRETURN a, b"
		if result.Query != expectedQuery {
			t.Errorf("Expected query '%s', but got '%s'", expectedQuery, result.Query)
		}
	})
}

func TestWhereClauseOperators(t *testing.T) {
	t.Run("WHERE with STARTS WITH, ENDS WITH, CONTAINS, and IN", func(t *testing.T) {
		qb := NewQueryBuilder()
		conditions := []types.Condition{
			types.Predicate{Property: "n.name", Operator: types.OpStartsWith, Value: "J"},
			types.Predicate{Property: "n.name", Operator: types.OpEndsWith, Value: "n"},
			types.Predicate{Property: "n.name", Operator: types.OpContains, Value: "oh"},
			types.Predicate{Property: "n.status", Operator: types.OpIn, Value: []string{"active", "pending"}},
		}

		result, err := qb.Match("(n:Person)").As("n").Where(conditions...).Return("n").Build()
		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		expectedQuery := "MATCH (n:Person)\nWHERE n.name STARTS WITH $n_name AND n.name ENDS WITH $n_name AND n.name CONTAINS $n_name AND n.status IN $n_status_list\nRETURN n"
		if result.Query != expectedQuery {
			t.Errorf("Expected query '%s', but got '%s'", expectedQuery, result.Query)
		}
	})
}

func TestUnwindClause(t *testing.T) {
	t.Run("UNWIND list", func(t *testing.T) {
		qb := NewQueryBuilder()
		result, err := qb.Unwind("[1, 2, 3]", "x").
			Return("x").
			Build()
		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		expectedQuery := "UNWIND [1, 2, 3] AS x\nRETURN x"
		if result.Query != expectedQuery {
			t.Errorf("Expected query '%s', but got '%s'", expectedQuery, result.Query)
		}
	})
}

func TestOrderBySkipLimitClauses(t *testing.T) {
	t.Run("ORDER BY, SKIP, and LIMIT", func(t *testing.T) {
		qb := NewQueryBuilder()
		result, err := qb.Match("(n:Person)").
			Return("n.name, n.age").
			OrderBy("n.age DESC", "n.name ASC").
			Skip(10).
			Limit(20).
			Build()
		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		expectedQuery := "MATCH (n:Person)\nRETURN n.name, n.age\nORDER BY n.age DESC, n.name ASC\nSKIP 10\nLIMIT 20"
		if result.Query != expectedQuery {
			t.Errorf("Expected query '%s', but got '%s'", expectedQuery, result.Query)
		}
	})
}



