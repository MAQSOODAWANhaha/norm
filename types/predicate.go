// types/predicate.go
package types

// Operator represents a comparison or logical operator.
type Operator string

const (
	// Comparison Operators
	OpEqual              Operator = "="
	OpNotEqual           Operator = "<>"
	OpLessThan           Operator = "<"
	OpLessThanOrEqual    Operator = "<="
	OpGreaterThan        Operator = ">"
	OpGreaterThanOrEqual Operator = ">="
	OpContains           Operator = "CONTAINS"
	OpStartsWith         Operator = "STARTS WITH"
	OpEndsWith           Operator = "ENDS WITH"
	OpRegex              Operator = "=~"
	OpIn                 Operator = "IN"
	OpIsNull             Operator = "IS NULL"
	OpIsNotNull          Operator = "IS NOT NULL"

	// Property Existence Operator
	OpExists Operator = "EXISTS"

	// Logical Operators
	OpAnd Operator = "AND"
	OpOr  Operator = "OR"
	OpNot Operator = "NOT"

	// Set Operator
	OpSet Operator = "+="
)

// Condition represents a part of a WHERE clause, which can be a simple
// predicate (e.g., name = 'Alice') or a logical grouping of other conditions.
type Condition interface {
	// isCondition is a private method to ensure only types from this package
	// can implement this interface.
	isCondition()
}

// Predicate is a basic condition, like "property operator value".
// e.g., "u.name = 'Alice'".
type Predicate struct {
	Property string
	Operator Operator
	Value    interface{} // Can be a single value or a slice for IN.
	Not      bool        // Used for prepending NOT to the predicate.
}

func (p Predicate) isCondition() {}

// LogicalGroup is a collection of conditions joined by a logical operator (AND/OR).
// e.g., "(u.age > 25 AND u.active = true)".
type LogicalGroup struct {
	Operator   Operator
	Conditions []Condition
}

func (lg LogicalGroup) isCondition() {}

// ExistsClause represents an EXISTS subquery.
// e.g., "EXISTS { MATCH (n)-[:KNOWS]->(m) }".
type ExistsClause struct {
	Query QueryBuilder
}

func (e ExistsClause) isCondition() {}

// QueryBuilder is an interface that represents a query builder.
// This is needed to avoid circular dependencies.
type QueryBuilder interface {
	Build() (QueryResult, error)
}