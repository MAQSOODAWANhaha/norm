// tests/query_builder_test.go
package tests

import (
	"strings"
	"testing"

	"norm/builder"
)

func TestQueryBuilder(t *testing.T) {
	t.Run("Basic Create Entity", func(t *testing.T) {
		user := &User{Username: "testuser", Email: "test@example.com", Active: true}
		result, err := builder.NewQueryBuilder().
			Create(user).As("u").
			Return("u").
			Build()

		if err != nil {
			t.Fatalf("Build() failed: %v", err)
		}

		// A more robust check for the query string
		if !strings.Contains(result.Query, "CREATE (u:User:Person") || !strings.Contains(result.Query, "RETURN u") {
			t.Errorf("unexpected query structure: %s", result.Query)
		}

		expectedParams := map[string]interface{}{
			"active_1":   true,
			"email_2":    "test@example.com",
			"username_3": "testuser",
		}

		if len(result.Parameters) != len(expectedParams) {
			t.Errorf("unexpected number of parameters: got %d, want %d", len(result.Parameters), len(expectedParams))
		}
	})

	t.Run("Create Entity with Default Label", func(t *testing.T) {
		type NoLabelUser struct {
			Name string `cypher:"name"`
		}
		user := &NoLabelUser{Name: "defaultUser"}
		result, err := builder.NewQueryBuilder().
			Create(user).As("u").
			Return("u").
			Build()

		if err != nil {
			t.Fatalf("Build() failed: %v", err)
		}

		if !strings.Contains(result.Query, "CREATE (u:NoLabelUser {name: $name_1})") || !strings.Contains(result.Query, "RETURN u") {
			t.Errorf("unexpected query structure for default label: %s", result.Query)
		}

		if result.Parameters["name_1"] != "defaultUser" {
			t.Errorf("unexpected parameter for default label user: got %v, want %s", result.Parameters["name_1"], "defaultUser")
		}
	})
}

func TestQueryBuilder_Set(t *testing.T) {
	t.Run("Simple Set", func(t *testing.T) {
		b := builder.NewQueryBuilder()
        result, err := b.
            Match(&User{}).As("u").
            Where(builder.Eq("u.username", "testuser")).
            Set(map[string]interface{}{"active": false}).
            Return("u").
            Build()

        if err != nil {
            t.Fatalf("Build() failed: %v", err)
        }

        expectedQuery := "MATCH (u:User:Person)\nWHERE (u.username = $u_username_1)\nSET u.active = $active_2\nRETURN u"
        if result.Query != expectedQuery {
            t.Errorf("unexpected query string:\ngot:  %s\nwant: %s", result.Query, expectedQuery)
        }

        expectedParams := map[string]interface{}{
            "u_username_1": "testuser",
            "active_2":     false,
        }

        if len(result.Parameters) != len(expectedParams) {
            t.Errorf("unexpected number of parameters: got %d, want %d", len(result.Parameters), len(expectedParams))
        }
        for k, v := range expectedParams {
            if result.Parameters[k] != v {
                t.Errorf("parameter %s mismatch: got %v, want %v", k, result.Parameters[k], v)
            }
        }
    })
}
