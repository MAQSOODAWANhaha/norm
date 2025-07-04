// validator/query.go
package validator

import (
	"strings"

	"norm/types"
)

// QueryValidator 查询验证器接口
type QueryValidator interface {
	Validate(query string) []types.ValidationError
	ValidateStructure(clauses []types.Clause) []types.ValidationError
	ValidateParameters(params map[string]interface{}) []types.ValidationError
}

// cypherQueryValidator 实现 QueryValidator 接口
type cypherQueryValidator struct {
	strictMode bool
}

// NewQueryValidator 创建新的查询验证器
func NewQueryValidator(strictMode bool) QueryValidator {
	return &cypherQueryValidator{
		strictMode: strictMode,
	}
}

// Validate 验证完整查询
func (v *cypherQueryValidator) Validate(query string) []types.ValidationError {
	var errors []types.ValidationError

	// 基本语法检查
	if query == "" {
		errors = append(errors, types.ValidationError{
			Type:       "empty_query",
			Message:    "Query cannot be empty",
			Position:   0,
			Suggestion: "Provide a valid Cypher query",
		})
		return errors
	}

	// 检查括号匹配
	if !v.validateBrackets(query) {
		errors = append(errors, types.ValidationError{
			Type:       "bracket_mismatch",
			Message:    "Mismatched brackets",
			Position:   -1, // Position is hard to determine accurately without a full parser
			Suggestion: "Check that all parentheses (), square brackets [], and curly braces {} are correctly paired",
		})
	}

	// 检查关键字使用
	errors = append(errors, v.validateKeywords(query)...)

	return errors
}

// validateBrackets 验证括号匹配
func (v *cypherQueryValidator) validateBrackets(query string) bool {
	stack := make([]rune, 0)
	pairs := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, char := range query {
		switch char {
		case '(', '[', '{':
			stack = append(stack, char)
		case ')', ']', '}':
			if len(stack) == 0 || stack[len(stack)-1] != pairs[char] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

// validateKeywords 验证关键字使用
func (v *cypherQueryValidator) validateKeywords(query string) []types.ValidationError {
	var errors []types.ValidationError

	// 检查是否有有效的子句
	validClauses := []string{"MATCH", "CREATE", "MERGE", "RETURN", "WITH", "WHERE", "UNWIND", "SET", "DELETE", "REMOVE"}
	hasValidClause := false

	upperQuery := strings.ToUpper(query)
	for _, clause := range validClauses {
		if strings.Contains(upperQuery, clause) {
			hasValidClause = true
			break
		}
	}

	if !hasValidClause {
		errors = append(errors, types.ValidationError{
			Type:       "no_valid_clause",
			Message:    "Query must contain at least one valid Cypher clause",
			Position:   0,
			Suggestion: "Add MATCH, CREATE, MERGE, or another valid clause",
		})
	}

	return errors
}

// ValidateStructure 验证子句结构 (暂未实现)
func (v *cypherQueryValidator) ValidateStructure(clauses []types.Clause) []types.ValidationError {
	// TODO: Implement structural validation, e.g., RETURN should be the last clause.
	return nil
}

// ValidateParameters 验证查询参数 (暂未实现)
func (v *cypherQueryValidator) ValidateParameters(params map[string]interface{}) []types.ValidationError {
	// TODO: Implement parameter validation.
	return nil
}