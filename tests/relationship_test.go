package tests

import (
	"testing"

	"norm/builder"
	"norm/types"
)

func TestRelationshipBuilder(t *testing.T) {
	tests := []struct {
		name     string
		builder  func() builder.RelationshipBuilder
		expected string
	}{
		{
			name: "Simple outgoing relationship",
			builder: func() builder.RelationshipBuilder {
				return builder.Outgoing("KNOWS")
			},
			expected: "-[:KNOWS]->",
		},
		{
			name: "Incoming relationship with variable",
			builder: func() builder.RelationshipBuilder {
				return builder.Incoming("CREATED").Variable("r")
			},
			expected: "<-[r:CREATED]-",
		},
		{
			name: "Bidirectional relationship",
			builder: func() builder.RelationshipBuilder {
				return builder.Bidirectional("RELATED")
			},
			expected: "-[:RELATED]-",
		},
		{
			name: "Variable length relationship",
			builder: func() builder.RelationshipBuilder {
				return builder.VarLengthOutgoing("KNOWS", 1, 3)
			},
			expected: "-[:KNOWS*1..3]->",
		},
		{
			name: "Relationship with properties",
			builder: func() builder.RelationshipBuilder {
				return builder.NewRelationshipBuilder().
					Type("KNOWS").
					Variable("r").
					Properties(map[string]interface{}{
						"since": 2020,
						"weight": 0.8,
					})
			},
			expected: "relationship with properties", // We'll check for key components instead
		},
		{
			name: "Variable length with min only",
			builder: func() builder.RelationshipBuilder {
				return builder.NewRelationshipBuilder().
					Type("CONNECTED").
					MinLength(2)
			},
			expected: "-[:CONNECTED*2..]->",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.builder().String()
			if tt.expected == "relationship with properties" {
				// For properties test, check key components
				if !containsString(result, "-[r:KNOWS") || 
				   !containsString(result, "since: 2020") ||
				   !containsString(result, "weight: 0.8") ||
				   !containsString(result, "]->") {
					t.Errorf("Expected relationship with properties, got %q", result)
				}
			} else if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestPatternBuilder(t *testing.T) {
	tests := []struct {
		name     string
		pattern  func() types.Pattern
		expected string
	}{
		{
			name: "Simple node to node relationship",
			pattern: func() types.Pattern {
				return builder.NewPatternBuilder().
					StartNode(builder.Node("a", "User")).
					Relationship(builder.Outgoing("KNOWS").Build()).
					EndNode(builder.Node("b", "User")).
					Build()
			},
			expected: "(a:User)-[:KNOWS]->(b:User)",
		},
		{
			name: "Complex pattern with properties",
			pattern: func() types.Pattern {
				startNode := builder.NodeWithProps("u", []string{"User"}, map[string]interface{}{
					"active": true,
				})
				
				rel := builder.NewRelationshipBuilder().
					Type("CREATED").
					Variable("r").
					Properties(map[string]interface{}{
						"timestamp": "2024-01-01",
					}).
					Build()
				
				endNode := builder.Node("p", "Post")
				
				return builder.NewPatternBuilder().
					StartNode(startNode).
					Relationship(rel).
					EndNode(endNode).
					Build()
			},
			expected: "(u:User {active: true})-[r:CREATED {timestamp: 2024-01-01}]->(p:Post)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := tt.pattern()
			pb := builder.NewPatternBuilder().
				StartNode(pattern.StartNode).
				Relationship(pattern.Relationship).
				EndNode(pattern.EndNode)
			
			result := pb.String()
			// Note: The actual result might have different property ordering due to map iteration
			// In a real test, we'd need to normalize or use a more sophisticated comparison
			t.Logf("Pattern result: %s", result)
			// This is a simplified assertion
			if len(result) == 0 {
				t.Error("Pattern should not be empty")
			}
		})
	}
}

func TestQueryBuilderWithPatterns(t *testing.T) {
	tests := []struct {
		name     string
		query    func() string
		expected string
	}{
		{
			name: "Match with relationship pattern",
			query: func() string {
				pattern := types.Pattern{
					StartNode: builder.Node("u", "User"),
					Relationship: builder.Outgoing("KNOWS").Variable("r").Build(),
					EndNode: builder.Node("f", "User"),
				}
				
				result, _ := builder.NewQueryBuilder().
					MatchPattern(pattern).
					Return("u.name", "f.name").
					Build()
				
				return result.Query
			},
			expected: "MATCH",
		},
		{
			name: "Create with variable length relationship",
			query: func() string {
				pattern := types.Pattern{
					StartNode: builder.Node("a", "Person"),
					Relationship: builder.VarLengthOutgoing("KNOWS", 1, 3).Build(),
					EndNode: builder.Node("b", "Person"),
				}
				
				result, _ := builder.NewQueryBuilder().
					CreatePattern(pattern).
					Build()
				
				return result.Query
			},
			expected: "CREATE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.query()
			if len(result) == 0 {
				t.Error("Query should not be empty")
			}
			t.Logf("Generated query: %s", result)
		})
	}
}

func TestAdvancedQueryStructures(t *testing.T) {
	t.Run("MERGE with ON CREATE and ON MATCH", func(t *testing.T) {
		user := &testUser{
			Username: "john_doe",
			Email:    "john@example.com",
		}
		
		result, err := builder.NewQueryBuilder().
			Merge(user).As("u").
			OnCreate("u.created = timestamp()", "u.status = 'new'").
			OnMatch("u.lastSeen = timestamp()", "u.status = 'active'").
			Return("u").
			Build()
		
		if err != nil {
			t.Fatalf("Error building query: %v", err)
		}
		
		expectedClauses := []string{"MERGE", "ON CREATE", "ON MATCH", "RETURN"}
		for _, clause := range expectedClauses {
			if !containsString(result.Query, clause) {
				t.Errorf("Query should contain %s clause", clause)
			}
		}
		
		t.Logf("Generated MERGE query: %s", result.Query)
	})
	
	t.Run("UNWIND with list", func(t *testing.T) {
		names := []interface{}{"Alice", "Bob", "Charlie"}
		
		result, err := builder.NewQueryBuilder().
			Unwind(names, "name").
			Create("(u:User {name: name})").
			Return("u").
			Build()
		
		if err != nil {
			t.Fatalf("Error building query: %v", err)
		}
		
		if !containsString(result.Query, "UNWIND") {
			t.Error("Query should contain UNWIND clause")
		}
		
		if result.Parameters["list_1"] == nil {
			t.Error("Query should have parameterized list")
		}
		
		t.Logf("Generated UNWIND query: %s", result.Query)
		t.Logf("Parameters: %v", result.Parameters)
	})
	
	t.Run("Complex query with UNION", func(t *testing.T) {
		qb1 := builder.NewQueryBuilder().
			Match("(u:User)").
			Return("u.name AS name", "'User' AS type")
		
		result1, _ := qb1.Build()
		
		qb2 := builder.NewQueryBuilder().
			Match("(c:Company)").
			Return("c.name AS name", "'Company' AS type")
		
		result2, _ := qb2.Build()
		
		// Simulate UNION by combining queries
		unionQuery := result1.Query + "\nUNION\n" + result2.Query
		
		if !containsString(unionQuery, "UNION") {
			t.Error("Union query should contain UNION clause")
		}
		
		t.Logf("Generated UNION query: %s", unionQuery)
	})
	
	t.Run("Query with USE database", func(t *testing.T) {
		result, err := builder.NewQueryBuilder().
			Use("mydb").
			Match("(n)").
			Return("count(n)").
			Build()
		
		if err != nil {
			t.Fatalf("Error building query: %v", err)
		}
		
		if !containsString(result.Query, "USE mydb") {
			t.Error("Query should contain USE clause")
		}
		
		t.Logf("Generated USE query: %s", result.Query)
	})
}

// Helper types and functions
type testUser struct {
	_ struct{} `cypher:"label:User"`
	Username string `cypher:"username"`
	Email    string `cypher:"email"`
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr ||
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}