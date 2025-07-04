// builder/types.go
package builder

// QueryType 表示不同类型的查询
type QueryType string

const (
	ReadQuery      QueryType = "read"
	WriteQuery     QueryType = "write"
	ReadWriteQuery QueryType = "read_write"
)

// Operator 表示查询操作符
type Operator string

const (
	// 比较操作符
	OpEqual              Operator = "="
	OpNotEqual           Operator = "<>"
	OpLessThan           Operator = "<"
	OpLessThanOrEqual    Operator = "<="
	OpGreaterThan        Operator = ">"
	OpGreaterThanOrEqual Operator = ">="

	// 字符串操作符
	OpContains   Operator = "CONTAINS"
	OpStartsWith Operator = "STARTS WITH"
	OpEndsWith   Operator = "ENDS WITH"
	OpRegex      Operator = "=~"

	// 列表操作符
	OpIn Operator = "IN"

	// 空值操作符
	OpIsNull    Operator = "IS NULL"
	OpIsNotNull Operator = "IS NOT NULL"

	// 逻辑操作符
	OpAnd Operator = "AND"
	OpOr  Operator = "OR"
	OpNot Operator = "NOT"
	OpXor Operator = "XOR"
)
