// builder/expression_test.go
package builder

import (
	"testing"
)

func TestPathFunctions(t *testing.T) {
	t.Run("shortestPath function", func(t *testing.T) {
		expr := ShortestPath("(a)-[:KNOWS*]-(b)")
		if expr.String() != "shortestPath((a)-[:KNOWS*]-(b))" {
			t.Errorf("Expected 'shortestPath((a)-[:KNOWS*]-(b))', but got '%s'", expr.String())
		}
	})

	t.Run("allShortestPaths function", func(t *testing.T) {
		expr := AllShortestPaths("(a)-[:KNOWS*]-(b)")
		if expr.String() != "allShortestPaths((a)-[:KNOWS*]-(b))" {
			t.Errorf("Expected 'allShortestPaths((a)-[:KNOWS*]-(b))', but got '%s'", expr.String())
		}
	})
}
