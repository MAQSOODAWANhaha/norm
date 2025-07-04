// builder/query.go
package builder

import (
	"fmt"
	"reflect"
	"strings"

	"norm/model"
	"norm/types"
	"norm/validator"
)

// QueryBuilder 查询构建器接口
type QueryBuilder interface {
	// 基本子句
	Match(pattern string) QueryBuilder
	OptionalMatch(pattern string) QueryBuilder
	Create(pattern string) QueryBuilder
	Merge(pattern string) QueryBuilder
	Where(condition string) QueryBuilder

	// 表达式支持子句
	Return(expressions ...interface{}) QueryBuilder
	With(expressions ...interface{}) QueryBuilder

	// 排序和限制
	OrderBy(fields ...string) QueryBuilder
	Skip(count int) QueryBuilder
	Limit(count int) QueryBuilder

	// 参数操作
	SetParameter(key string, value interface{}) QueryBuilder

	// 构建操作
	Build() (types.QueryResult, error)
	Validate() []types.ValidationError

	// 实体操作
	MatchEntity(entity interface{}) QueryBuilder
	CreateEntity(entity interface{}) QueryBuilder
}

// cypherQueryBuilder 实现 QueryBuilder 接口
type cypherQueryBuilder struct {
	clauses      []types.Clause
	parameters   map[string]interface{}
	paramCounter int
	registry     model.Registry
	validator    validator.QueryValidator
}

// NewQueryBuilder 创建新的查询构建器
func NewQueryBuilder(registry model.Registry) QueryBuilder {
	return &cypherQueryBuilder{
		clauses:      make([]types.Clause, 0),
		parameters:   make(map[string]interface{}),
		paramCounter: 0,
		registry:     registry,
		validator:    validator.NewQueryValidator(true), // strictMode = true
	}
}

func (q *cypherQueryBuilder) addClause(clauseType types.ClauseType, content string) QueryBuilder {
	q.clauses = append(q.clauses, types.Clause{
		Type:    clauseType,
		Content: content,
	})
	return q
}

// Match 添加 MATCH 子句
func (q *cypherQueryBuilder) Match(pattern string) QueryBuilder {
	return q.addClause(types.MatchClause, pattern)
}

// OptionalMatch 添加 OPTIONAL MATCH 子句
func (q *cypherQueryBuilder) OptionalMatch(pattern string) QueryBuilder {
	return q.addClause(types.OptionalMatchClause, pattern)
}

// Create 添加 CREATE 子句
func (q *cypherQueryBuilder) Create(pattern string) QueryBuilder {
	return q.addClause(types.CreateClause, pattern)
}

// Merge 添加 MERGE 子句
func (q *cypherQueryBuilder) Merge(pattern string) QueryBuilder {
	return q.addClause(types.MergeClause, pattern)
}

// Where 添加 WHERE 子句
func (q *cypherQueryBuilder) Where(condition string) QueryBuilder {
	return q.addClause(types.WhereClause, condition)
}

// Return 添加 RETURN 子句
func (q *cypherQueryBuilder) Return(expressions ...interface{}) QueryBuilder {
	return q.addClause(types.ReturnClause, q.formatExpressions(expressions...))
}

// With 添加 WITH 子句
func (q *cypherQueryBuilder) With(expressions ...interface{}) QueryBuilder {
	return q.addClause(types.WithClause, q.formatExpressions(expressions...))
}

// OrderBy 添加 ORDER BY 子句
func (q *cypherQueryBuilder) OrderBy(fields ...string) QueryBuilder {
	return q.addClause(types.OrderByClause, strings.Join(fields, ", "))
}

// Skip 添加 SKIP 子句
func (q *cypherQueryBuilder) Skip(count int) QueryBuilder {
	return q.addClause(types.SkipClause, fmt.Sprintf("%d", count))
}

// Limit 添加 LIMIT 子句
func (q *cypherQueryBuilder) Limit(count int) QueryBuilder {
	return q.addClause(types.LimitClause, fmt.Sprintf("%d", count))
}

// SetParameter 设置查询参数
func (q *cypherQueryBuilder) SetParameter(key string, value interface{}) QueryBuilder {
	q.parameters[key] = value
	return q
}

// MatchEntity 匹配实体
func (q *cypherQueryBuilder) MatchEntity(entity interface{}) QueryBuilder {
	pattern, err := q.buildEntityPattern(entity, "")
	if err != nil {
		// 错误处理将在第二阶段完善
		return q
	}
	return q.Match(pattern)
}

// CreateEntity 创建实体
func (q *cypherQueryBuilder) CreateEntity(entity interface{}) QueryBuilder {
	pattern, err := q.buildEntityPattern(entity, "")
	if err != nil {
		// 错误处理将在第二阶段完善
		return q
	}
	return q.Create(pattern)
}

// buildEntityPattern 构建实体模式
func (q *cypherQueryBuilder) buildEntityPattern(entity interface{}, variable string) (string, error) {
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	metadata, exists := q.registry.GetByType(t)
	if !exists {
		return "", fmt.Errorf("entity %s not registered", t.Name())
	}

	nodeBuilder := NewNodeBuilder()
	if variable != "" {
		nodeBuilder = nodeBuilder.Variable(variable)
	}
	nodeBuilder = nodeBuilder.Labels(metadata.Labels...)

	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	properties := make(map[string]interface{})
	for fieldName, propMeta := range metadata.Properties {
		fieldValue := v.FieldByName(fieldName)
		if fieldValue.IsValid() && !fieldValue.IsZero() {
			paramName := q.generateParameterName(propMeta.CypherName)
			properties[propMeta.CypherName] = fmt.Sprintf("$"+"%s", paramName)
			q.parameters[paramName] = fieldValue.Interface()
		}
	}

	if len(properties) > 0 {
		nodeBuilder = nodeBuilder.Properties(properties)
	}

	return nodeBuilder.Build(), nil
}

// generateParameterName 生成参数名
func (q *cypherQueryBuilder) generateParameterName(base string) string {
	q.paramCounter++
	return fmt.Sprintf("%s_%d", base, q.paramCounter)
}

// formatExpressions 格式化表达式
func (q *cypherQueryBuilder) formatExpressions(expressions ...interface{}) string {
	var parts []string
	for _, expr := range expressions {
		switch v := expr.(type) {
		case string:
			parts = append(parts, v)
		case Expression:
			parts = append(parts, v.String())
		default:
			parts = append(parts, fmt.Sprintf("%v", v))
		}
	}
	return strings.Join(parts, ", ")
}

// Build 构建最终的 Cypher 查询
func (q *cypherQueryBuilder) Build() (types.QueryResult, error) {
	var parts []string
	for _, clause := range q.clauses {
		part := string(clause.Type)
		if clause.Content != "" {
			part += " " + clause.Content
		}
		parts = append(parts, part)
	}

	query := strings.Join(parts, "\n")
	errors := q.Validate()

	return types.QueryResult{
		Query:      query,
		Parameters: q.parameters,
		Valid:      len(errors) == 0,
		Errors:     errors,
	}, nil
}

// Validate 验证查询
func (q *cypherQueryBuilder) Validate() []types.ValidationError {
	var parts []string
	for _, clause := range q.clauses {
		part := string(clause.Type)
		if clause.Content != "" {
			part += " " + clause.Content
		}
		parts = append(parts, part)
	}
	query := strings.Join(parts, "\n")

	return q.validator.Validate(query)
}
