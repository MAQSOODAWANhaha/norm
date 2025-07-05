// builder/query.go
package builder

import (
	"fmt"
	"sort"
	"strings"

	"norm/types"
	"norm/validator"
)

// QueryBuilder is the interface for the Cypher query builder.
type QueryBuilder interface {
	// 基本模式匹配
	Match(patternOrEntity interface{}) QueryBuilder
	OptionalMatch(patternOrEntity interface{}) QueryBuilder
	Create(patternOrEntity interface{}) QueryBuilder
	Merge(patternOrEntity interface{}) QueryBuilder
	As(alias string) QueryBuilder
	
	// 关系模式支持
	MatchPattern(pattern types.Pattern) QueryBuilder
	CreatePattern(pattern types.Pattern) QueryBuilder
	MergePattern(pattern types.Pattern) QueryBuilder
	
	// 数据修改
	Set(assignments ...string) QueryBuilder
	SetEntity(entity interface{}, alias string) QueryBuilder
	Delete(variables ...interface{}) QueryBuilder
	DetachDelete(variables ...interface{}) QueryBuilder
	Remove(items ...string) QueryBuilder
	RemoveProperties(entity interface{}, alias string, properties ...string) QueryBuilder
	
	// MERGE 条件动作
	OnCreate(assignments ...string) QueryBuilder
	OnMatch(assignments ...string) QueryBuilder
	
	// 条件和过滤
	Where(conditions ...types.Condition) QueryBuilder
	WhereString(condition string) QueryBuilder
	
	// 数据返回和处理
	Return(expressions ...interface{}) QueryBuilder
	With(expressions ...interface{}) QueryBuilder
	Unwind(list interface{}, alias string) QueryBuilder
	
	// 排序和限制
	OrderBy(fields ...string) QueryBuilder
	Skip(count int) QueryBuilder
	Limit(count int) QueryBuilder
	
	// 集合操作
	Union() QueryBuilder
	UnionAll() QueryBuilder
	
	// 高级功能
	Use(database string) QueryBuilder
	Call(subquery QueryBuilder) QueryBuilder
	ForEach(variable string, list interface{}, updateClauses ...string) QueryBuilder
	
	// 参数和构建
	SetParameter(key string, value interface{}) QueryBuilder
	Build() (types.QueryResult, error)
	Validate() []types.ValidationError
}

// cypherQueryBuilder implements the QueryBuilder interface.
type cypherQueryBuilder struct {
	clauses       []types.Clause
	parameters    map[string]interface{}
	paramCounter  int
	currentAlias  string
	pendingEntity interface{}
	pendingClause types.ClauseType
	entityAliases map[string]interface{}
	validator     validator.QueryValidator
	errors        []error
}

// NewQueryBuilder creates a new instance of the query builder.
func NewQueryBuilder() QueryBuilder {
	return &cypherQueryBuilder{
		clauses:       make([]types.Clause, 0),
		parameters:    make(map[string]interface{}),
		paramCounter:  0,
		entityAliases: make(map[string]interface{}),
		validator:     validator.NewQueryValidator(true),
		errors:        make([]error, 0),
	}
}

// handleEntityClause handles methods that can take a string pattern or an entity struct.
func (q *cypherQueryBuilder) handleEntityClause(clauseType types.ClauseType, p interface{}) QueryBuilder {
	q.finalizePendingClause()
	switch v := p.(type) {
	case string:
		q.addClause(clauseType, v)
	default:
		q.pendingEntity = v
		q.pendingClause = clauseType
	}
	return q
}

func (q *cypherQueryBuilder) Match(p interface{}) QueryBuilder {
	return q.handleEntityClause(types.MatchClause, p)
}

func (q *cypherQueryBuilder) OptionalMatch(p interface{}) QueryBuilder {
	return q.handleEntityClause(types.OptionalMatchClause, p)
}

func (q *cypherQueryBuilder) Create(p interface{}) QueryBuilder {
	return q.handleEntityClause(types.CreateClause, p)
}

func (q *cypherQueryBuilder) Merge(p interface{}) QueryBuilder {
	return q.handleEntityClause(types.MergeClause, p)
}

// As sets the alias for a pending entity clause.
func (q *cypherQueryBuilder) As(alias string) QueryBuilder {
	if q.pendingEntity == nil {
		q.currentAlias = alias
		return q
	}
	q.currentAlias = alias
	q.entityAliases[alias] = q.pendingEntity
	q.finalizePendingClause()
	return q
}

func (q *cypherQueryBuilder) Set(assignments ...string) QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.SetClause, strings.Join(assignments, ", "))
	return q
}

func (q *cypherQueryBuilder) SetEntity(entity interface{}, alias string) QueryBuilder {
	q.finalizePendingClause()
	props, err := ParseEntityForUpdate(entity)
	if err != nil {
		q.errors = append(q.errors, err)
		return q
	}

	var assignments []string
	for key, value := range props {
		paramName := q.generateParameterName(key)
		assignments = append(assignments, fmt.Sprintf("%s.%s = $%s", alias, key, paramName))
		q.parameters[paramName] = value
	}

	if len(assignments) > 0 {
		q.addClause(types.SetClause, strings.Join(assignments, ", "))
	}
	return q
}

// finalizePendingClause builds and adds the clause that was waiting for an alias.
func (q *cypherQueryBuilder) finalizePendingClause() {
	if q.pendingEntity == nil {
		return
	}

	pattern, err := q.buildEntityPattern(q.pendingEntity, q.currentAlias, q.pendingClause)
	if err != nil {
		q.errors = append(q.errors, err)
	} else {
		q.addClause(q.pendingClause, pattern)
	}

	q.pendingEntity = nil
	q.pendingClause = ""
}

func (q *cypherQueryBuilder) Where(conditions ...types.Condition) QueryBuilder {
	q.finalizePendingClause()
	if len(conditions) == 0 {
		return q
	}
	
	var conditionStr strings.Builder
	
	// Always create an AND group for consistent formatting
	group := types.LogicalGroup{Operator: types.OpAnd, Conditions: conditions}
	q.buildConditionString(&group, &conditionStr)

	q.addClause(types.WhereClause, conditionStr.String())
	return q
}

func (q *cypherQueryBuilder) WhereString(condition string) QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.WhereClause, condition)
	return q
}

func (q *cypherQueryBuilder) Return(expressions ...interface{}) QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.ReturnClause, q.formatExpressions(expressions...))
	return q
}

func (q *cypherQueryBuilder) With(expressions ...interface{}) QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.WithClause, q.formatExpressions(expressions...))
	return q
}

func (q *cypherQueryBuilder) OrderBy(fields ...string) QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.OrderByClause, strings.Join(fields, ", "))
	return q
}

func (q *cypherQueryBuilder) Skip(count int) QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.SkipClause, fmt.Sprintf("%d", count))
	return q
}

func (q *cypherQueryBuilder) Limit(count int) QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.LimitClause, fmt.Sprintf("%d", count))
	return q
}

func (q *cypherQueryBuilder) SetParameter(key string, value interface{}) QueryBuilder {
	q.parameters[key] = value
	return q
}

func (q *cypherQueryBuilder) Call(subquery QueryBuilder) QueryBuilder {
	q.finalizePendingClause()

	sub, ok := subquery.(*cypherQueryBuilder)
	if !ok {
		q.errors = append(q.errors, fmt.Errorf("subquery is not a valid *cypherQueryBuilder"))
		return q
	}

	// The subquery should be built with its own context.
	// We pass the parameter counter to avoid name collisions.
	sub.paramCounter = q.paramCounter
	subResult, err := sub.Build()
	if err != nil {
		q.errors = append(q.errors, fmt.Errorf("failed to build subquery: %w", err))
		return q
	}
	q.paramCounter = sub.paramCounter

	// Merge parameters.
	for k, v := range subResult.Parameters {
		q.parameters[k] = v
	}

	// Add CALL clause with the subquery string.
	q.addClause(types.CallClause, fmt.Sprintf("{\n%s\n}", subResult.Query))

	return q
}

// 关系模式支持方法
func (q *cypherQueryBuilder) MatchPattern(pattern types.Pattern) QueryBuilder {
	q.finalizePendingClause()
	patternStr := q.buildPatternString(pattern)
	q.addClause(types.MatchClause, patternStr)
	return q
}

func (q *cypherQueryBuilder) CreatePattern(pattern types.Pattern) QueryBuilder {
	q.finalizePendingClause()
	patternStr := q.buildPatternString(pattern)
	q.addClause(types.CreateClause, patternStr)
	return q
}

func (q *cypherQueryBuilder) MergePattern(pattern types.Pattern) QueryBuilder {
	q.finalizePendingClause()
	patternStr := q.buildPatternString(pattern)
	q.addClause(types.MergeClause, patternStr)
	return q
}

// 数据修改方法
func (q *cypherQueryBuilder) Delete(variables ...interface{}) QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.DeleteClause, q.formatDeleteVariables(variables...))
	return q
}

func (q *cypherQueryBuilder) DetachDelete(variables ...interface{}) QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.DetachDeleteClause, q.formatDeleteVariables(variables...))
	return q
}

func (q *cypherQueryBuilder) Remove(items ...string) QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.RemoveClause, strings.Join(items, ", "))
	return q
}

func (q *cypherQueryBuilder) RemoveProperties(entity interface{}, alias string, properties ...string) QueryBuilder {
	q.finalizePendingClause()
	var itemsToRemove []string
	if len(properties) == 0 {
		// Remove all properties from the entity
		props, err := ParseEntityForReturn(entity, "")
		if err != nil {
			q.errors = append(q.errors, err)
			return q
		}
		for _, prop := range props {
			itemsToRemove = append(itemsToRemove, fmt.Sprintf("%s.%s", alias, prop))
		}
	} else {
		// Remove specified properties
		for _, prop := range properties {
			itemsToRemove = append(itemsToRemove, fmt.Sprintf("%s.%s", alias, prop))
		}
	}

	if len(itemsToRemove) > 0 {
		q.addClause(types.RemoveClause, strings.Join(itemsToRemove, ", "))
	}
	return q
}


// MERGE 条件动作方法
func (q *cypherQueryBuilder) OnCreate(assignments ...string) QueryBuilder {
	q.finalizePendingClause()
	if len(assignments) > 0 {
		content := fmt.Sprintf("SET %s", strings.Join(assignments, ", "))
		q.addClause(types.OnCreateClause, content)
	}
	return q
}

func (q *cypherQueryBuilder) OnMatch(assignments ...string) QueryBuilder {
	q.finalizePendingClause()
	if len(assignments) > 0 {
		content := fmt.Sprintf("SET %s", strings.Join(assignments, ", "))
		q.addClause(types.OnMatchClause, content)
	}
	return q
}

// 数据处理方法
func (q *cypherQueryBuilder) Unwind(list interface{}, alias string) QueryBuilder {
	q.finalizePendingClause()
	var listStr string
	switch v := list.(type) {
	case string:
		listStr = v
	case []interface{}:
		// 处理数组
		paramName := q.generateParameterName("list")
		q.parameters[paramName] = v
		listStr = fmt.Sprintf("$%s", paramName)
	default:
		listStr = fmt.Sprintf("%v", v)
	}
	q.addClause(types.UnwindClause, fmt.Sprintf("%s AS %s", listStr, alias))
	return q
}

// 集合操作方法
func (q *cypherQueryBuilder) Union() QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.UnionClause, "")
	return q
}

func (q *cypherQueryBuilder) UnionAll() QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.UnionAllClause, "")
	return q
}

// 高级功能方法
func (q *cypherQueryBuilder) Use(database string) QueryBuilder {
	q.finalizePendingClause()
	q.addClause(types.UseClause, database)
	return q
}

func (q *cypherQueryBuilder) ForEach(variable string, list interface{}, updateClauses ...string) QueryBuilder {
	q.finalizePendingClause()
	var listStr string
	switch v := list.(type) {
	case string:
		listStr = v
	case []interface{}:
		paramName := q.generateParameterName("foreach_list")
		q.parameters[paramName] = v
		listStr = fmt.Sprintf("$%s", paramName)
	default:
		listStr = fmt.Sprintf("%v", v)
	}
	
	clauseContent := fmt.Sprintf("(%s IN %s | %s)", variable, listStr, strings.Join(updateClauses, " "))
	q.addClause(types.ForEachClause, clauseContent)
	return q
}

func (q *cypherQueryBuilder) Build() (types.QueryResult, error) {
	q.finalizePendingClause()
	if len(q.errors) > 0 {
		// Join all errors into one
		var errStrings []string
		for _, err := range q.errors {
			errStrings = append(errStrings, err.Error())
		}
		return types.QueryResult{}, fmt.Errorf("%s", strings.Join(errStrings, "; "))
	}

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

// --- Helper Methods ---

func (q *cypherQueryBuilder) addClause(clauseType types.ClauseType, content string) {
	q.clauses = append(q.clauses, types.Clause{
		Type:    clauseType,
		Content: content,
	})
}

func (q *cypherQueryBuilder) buildEntityPattern(entity interface{}, variable string, clauseType types.ClauseType) (string, error) {
	entityInfo, err := ParseEntity(entity)
	if err != nil {
		return "", fmt.Errorf("failed to parse entity: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("(")
	if variable != "" {
		sb.WriteString(variable)
	}
	for _, label := range entityInfo.Labels {
		sb.WriteString(":")
		sb.WriteString(string(label))
	}

	// Only add properties for CREATE and MERGE clauses
	if (clauseType == types.CreateClause || clauseType == types.MergeClause) && len(entityInfo.Properties) > 0 {
		sb.WriteString(" {")
		var props []string
		
		// Sort keys for deterministic order
		keys := make([]string, 0, len(entityInfo.Properties))
		for k := range entityInfo.Properties {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		
		for _, k := range keys {
			paramName := q.generateParameterName(k)
			props = append(props, fmt.Sprintf("%s: $%s", k, paramName))
			q.parameters[paramName] = entityInfo.Properties[k]
		}
		sb.WriteString(strings.Join(props, ", "))
		sb.WriteString("}")
	}

	sb.WriteString(")")
	return sb.String(), nil
}

func (q *cypherQueryBuilder) buildConditionString(condition types.Condition, sb *strings.Builder) {
	switch c := condition.(type) {
	case types.Predicate:
		if c.Not {
			sb.WriteString("NOT (")
		}

		prop := c.Property
		// Don't modify property if it already contains a dot (already qualified)
		// Only add current alias if property doesn't contain dot and we have a current alias
		if !strings.Contains(prop, ".") && q.currentAlias != "" {
			prop = fmt.Sprintf("%s.%s", q.currentAlias, prop)
		}

		if c.Operator == types.OpIsNull || c.Operator == types.OpIsNotNull {
			sb.WriteString(fmt.Sprintf("%s %s", prop, c.Operator))
		} else {
			// Generate parameter name based on the full property (including alias if present)
			paramName := q.generateParameterName(strings.ReplaceAll(prop, ".", "_"))
			q.parameters[paramName] = c.Value
			sb.WriteString(fmt.Sprintf("%s %s $%s", prop, c.Operator, paramName))
		}

		if c.Not {
			sb.WriteString(")")
		}

	case types.LogicalGroup:
		sb.WriteString("(")
		for i, cond := range c.Conditions {
			if i > 0 {
				sb.WriteString(fmt.Sprintf(" %s ", c.Operator))
			}
			q.buildConditionString(cond, sb)
		}
		sb.WriteString(")")
	case *types.LogicalGroup:
		sb.WriteString("(")
		for i, cond := range c.Conditions {
			if i > 0 {
				sb.WriteString(fmt.Sprintf(" %s ", c.Operator))
			}
			q.buildConditionString(cond, sb)
		}
		sb.WriteString(")")
	}
}

func (q *cypherQueryBuilder) generateParameterName(base string) string {
	q.paramCounter++
	return fmt.Sprintf("%s_%d", strings.ReplaceAll(base, ".", "_"), q.paramCounter)
}

func (q *cypherQueryBuilder) formatExpressions(expressions ...interface{}) string {
	var parts []string
	for _, expr := range expressions {
		switch v := expr.(type) {
		case string:
			parts = append(parts, v)
		case Expression:
			parts = append(parts, v.String())
		case types.Entity:
			props, err := ParseEntityForReturn(v.Struct, v.Alias)
			if err != nil {
				q.errors = append(q.errors, err)
				continue
			}
			parts = append(parts, props...)
		default:
			parts = append(parts, fmt.Sprintf("%v", v))
		}
	}
	return strings.Join(parts, ", ")
}

func (q *cypherQueryBuilder) formatDeleteVariables(variables ...interface{}) string {
	var parts []string
	for _, v := range variables {
		switch val := v.(type) {
		case string:
			parts = append(parts, val)
		case types.Entity:
			parts = append(parts, val.Alias)
		default:
			// Attempt to find alias by struct type if not explicitly provided
			found := false
			for alias, entity := range q.entityAliases {
				if entity == val {
					parts = append(parts, alias)
					found = true
					break
				}
			}
			if !found {
				q.errors = append(q.errors, fmt.Errorf("could not find alias for entity to delete: %T", val))
			}
		}
	}
	return strings.Join(parts, ", ")
}


func (q *cypherQueryBuilder) buildPatternString(pattern types.Pattern) string {
	var sb strings.Builder
	
	// 起始节点
	sb.WriteString(q.buildNodePatternString(pattern.StartNode))
	
	// 关系
	sb.WriteString(q.buildRelationshipPatternString(pattern.Relationship))
	
	// 结束节点
	sb.WriteString(q.buildNodePatternString(pattern.EndNode))
	
	return sb.String()
}

func (q *cypherQueryBuilder) buildNodePatternString(node types.NodePattern) string {
	var sb strings.Builder
	sb.WriteString("(")
	
	// 变量名
	if node.Variable != "" {
		sb.WriteString(node.Variable)
	}
	
	// 标签
	for _, label := range node.Labels {
		sb.WriteString(":")
		sb.WriteString(string(label))
	}
	
	// 属性
	if len(node.Properties) > 0 {
		sb.WriteString(" {")
		var props []string
		
		// 排序属性键以确保确定性
		keys := make([]string, 0, len(node.Properties))
		for k := range node.Properties {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		
		for _, k := range keys {
			paramName := q.generateParameterName(k)
			props = append(props, fmt.Sprintf("%s: $%s", k, paramName))
			q.parameters[paramName] = node.Properties[k]
		}
		sb.WriteString(strings.Join(props, ", "))
		sb.WriteString("}")
	}
	
	sb.WriteString(")")
	return sb.String()
}

func (q *cypherQueryBuilder) buildRelationshipPatternString(rel types.RelationshipPattern) string {
	var sb strings.Builder
	
	// 开始方向
	switch rel.Direction {
	case types.DirectionIncoming:
		sb.WriteString("<-<")
	case types.DirectionOutgoing:
		sb.WriteString("-")
	case types.DirectionBoth:
		sb.WriteString("-")
	default:
		sb.WriteString("-")
	}
	
	sb.WriteString("[")
	
	// 变量名
	if rel.Variable != "" {
		sb.WriteString(rel.Variable)
	}
	
	// 关系类型
	if rel.Type != "" {
		sb.WriteString(":")
		sb.WriteString(rel.Type)
	}
	
	// 变长路径
	if rel.MinLength != nil || rel.MaxLength != nil {
		sb.WriteString("*")
		if rel.MinLength != nil {
			sb.WriteString(fmt.Sprintf("%d", *rel.MinLength))
		}
		if rel.MaxLength != nil {
			sb.WriteString("..<")
			sb.WriteString(fmt.Sprintf("%d", *rel.MaxLength))
		} else if rel.MinLength != nil {
			sb.WriteString("..<")
		}
	}
	
	// 属性
	if len(rel.Properties) > 0 {
		sb.WriteString(" {")
		var props []string
		
		// 排序属性键
		keys := make([]string, 0, len(rel.Properties))
		for k := range rel.Properties {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		
		for _, k := range keys {
			paramName := q.generateParameterName(k)
			props = append(props, fmt.Sprintf("%s: $%s", k, paramName))
			q.parameters[paramName] = rel.Properties[k]
		}
		sb.WriteString(strings.Join(props, ", "))
		sb.WriteString("}")
	}
	
	sb.WriteString("]")
	
	// 结束方向
	switch rel.Direction {
	case types.DirectionIncoming:
		sb.WriteString("-")
	case types.DirectionOutgoing:
		sb.WriteString("->")
	case types.DirectionBoth:
		sb.WriteString("-")
	default:
		sb.WriteString("->")
	}
	
	return sb.String()
}