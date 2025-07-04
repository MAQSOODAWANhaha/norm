// builder/query.go
package builder

import (
    "fmt"
    "reflect"
    "strings"
    "norm/types"
    "norm/model"
)

// QueryBuilder 查询构建器接口
type QueryBuilder interface {
    // 基本子句
    Match(pattern string) QueryBuilder
    OptionalMatch(pattern string) QueryBuilder
    Create(pattern string) QueryBuilder
    Merge(pattern string) QueryBuilder
    Set(assignments ...string) QueryBuilder
    Delete(nodes ...string) QueryBuilder
    Remove(properties ...string) QueryBuilder
    Return(fields ...string) QueryBuilder
    With(fields ...string) QueryBuilder
    Where(condition string) QueryBuilder
    
    // 排序和限制
    OrderBy(fields ...string) QueryBuilder
    OrderByDesc(fields ...string) QueryBuilder
    Skip(count int) QueryBuilder
    Limit(count int) QueryBuilder
    
    // 高级操作
    Unwind(expression, alias string) QueryBuilder
    Call(procedure string, params ...interface{}) QueryBuilder
    Union(other QueryBuilder) QueryBuilder
    UnionAll(other QueryBuilder) QueryBuilder
    
    // 参数操作
    SetParameter(key string, value interface{}) QueryBuilder
    SetParameters(params map[string]interface{}) QueryBuilder
    
    // 构建操作
    Build() (types.QueryResult, error)
    String() string
    
    // 实体操作
    MatchEntity(entity interface{}) QueryBuilder
    CreateEntity(entity interface{}) QueryBuilder
    MergeEntity(entity interface{}) QueryBuilder
    
    // 关系操作
    MatchRelationship(from, to interface{}, relType string) QueryBuilder
    CreateRelationship(from, to interface{}, relType string, props map[string]interface{}) QueryBuilder
}

// cypherQueryBuilder 实现 QueryBuilder 接口
type cypherQueryBuilder struct {
    clauses       []clause
    parameters    map[string]interface{}
    paramCounter  int
    registry      *model.EntityRegistry
    labelManager  *model.LabelManager
}

// clause 表示单个 Cypher 子句
type clause struct {
    Type      ClauseType
    Content   string
    Modifiers map[string]interface{}
}

// NewQueryBuilder 创建新的查询构建器
func NewQueryBuilder(registry *model.EntityRegistry) QueryBuilder {
    return &cypherQueryBuilder{
        clauses:      make([]clause, 0),
        parameters:   make(map[string]interface{}),
        paramCounter: 0,
        registry:     registry,
        labelManager: model.NewLabelManager(registry),
    }
}

// Match 添加 MATCH 子句
func (q *cypherQueryBuilder) Match(pattern string) QueryBuilder {
    q.clauses = append(q.clauses, clause{
        Type:    MatchClause,
        Content: pattern,
    })
    return q
}

// OptionalMatch 添加 OPTIONAL MATCH 子句
func (q *cypherQueryBuilder) OptionalMatch(pattern string) QueryBuilder {
    q.clauses = append(q.clauses, clause{
        Type:    OptionalMatchClause,
        Content: pattern,
    })
    return q
}

// Create 添加 CREATE 子句
func (q *cypherQueryBuilder) Create(pattern string) QueryBuilder {
    q.clauses = append(q.clauses, clause{
        Type:    CreateClause,
        Content: pattern,
    })
    return q
}

// Merge 添加 MERGE 子句
func (q *cypherQueryBuilder) Merge(pattern string) QueryBuilder {
    q.clauses = append(q.clauses, clause{
        Type:    MergeClause,
        Content: pattern,
    })
    return q
}

// Where 添加 WHERE 子句
func (q *cypherQueryBuilder) Where(condition string) QueryBuilder {
    q.clauses = append(q.clauses, clause{
        Type:    WhereClause,
        Content: condition,
    })
    return q
}

// Return 添加 RETURN 子句
func (q *cypherQueryBuilder) Return(fields ...string) QueryBuilder {
    content := strings.Join(fields, ", ")
    q.clauses = append(q.clauses, clause{
        Type:    ReturnClause,
        Content: content,
    })
    return q
}

// With 添加 WITH 子句
func (q *cypherQueryBuilder) With(fields ...string) QueryBuilder {
    content := strings.Join(fields, ", ")
    q.clauses = append(q.clauses, clause{
        Type:    WithClause,
        Content: content,
    })
    return q
}

// Set 添加 SET 子句
func (q *cypherQueryBuilder) Set(assignments ...string) QueryBuilder {
    content := strings.Join(assignments, ", ")
    q.clauses = append(q.clauses, clause{
        Type:    SetClause,
        Content: content,
    })
    return q
}

// Delete 添加 DELETE 子句
func (q *cypherQueryBuilder) Delete(nodes ...string) QueryBuilder {
    content := strings.Join(nodes, ", ")
    q.clauses = append(q.clauses, clause{
        Type:    DeleteClause,
        Content: content,
    })
    return q
}

// Remove 添加 REMOVE 子句
func (q *cypherQueryBuilder) Remove(properties ...string) QueryBuilder {
    content := strings.Join(properties, ", ")
    q.clauses = append(q.clauses, clause{
        Type:    RemoveClause,
        Content: content,
    })
    return q
}

// OrderBy 添加 ORDER BY 子句
func (q *cypherQueryBuilder) OrderBy(fields ...string) QueryBuilder {
    content := strings.Join(fields, ", ")
    q.clauses = append(q.clauses, clause{
        Type:    OrderByClause,
        Content: content,
    })
    return q
}

// OrderByDesc 添加 ORDER BY DESC 子句
func (q *cypherQueryBuilder) OrderByDesc(fields ...string) QueryBuilder {
    fieldList := make([]string, len(fields))
    for i, field := range fields {
        fieldList[i] = field + " DESC"
    }
    content := strings.Join(fieldList, ", ")
    q.clauses = append(q.clauses, clause{
        Type:    OrderByClause,
        Content: content,
    })
    return q
}

// Skip 添加 SKIP 子句
func (q *cypherQueryBuilder) Skip(count int) QueryBuilder {
    q.clauses = append(q.clauses, clause{
        Type:    SkipClause,
        Content: fmt.Sprintf("%d", count),
    })
    return q
}

// Limit 添加 LIMIT 子句
func (q *cypherQueryBuilder) Limit(count int) QueryBuilder {
    q.clauses = append(q.clauses, clause{
        Type:    LimitClause,
        Content: fmt.Sprintf("%d", count),
    })
    return q
}

// Unwind 添加 UNWIND 子句
func (q *cypherQueryBuilder) Unwind(expression, alias string) QueryBuilder {
    content := fmt.Sprintf("%s AS %s", expression, alias)
    q.clauses = append(q.clauses, clause{
        Type:    UnwindClause,
        Content: content,
    })
    return q
}

// Call 添加 CALL 子句
func (q *cypherQueryBuilder) Call(procedure string, params ...interface{}) QueryBuilder {
    content := procedure
    if len(params) > 0 {
        paramStrs := make([]string, len(params))
        for i, param := range params {
            paramName := q.generateParameterName()
            paramStrs[i] = "$" + paramName
            q.parameters[paramName] = param
        }
        content += "(" + strings.Join(paramStrs, ", ") + ")"
    }
    
    q.clauses = append(q.clauses, clause{
        Type:    CallClause,
        Content: content,
    })
    return q
}

// Union 添加 UNION 子句
func (q *cypherQueryBuilder) Union(other QueryBuilder) QueryBuilder {
    otherResult, _ := other.Build()
    q.clauses = append(q.clauses, clause{
        Type:    UnionClause,
        Content: "",
    })
    
    // 添加其他查询的子句
    q.clauses = append(q.clauses, clause{
        Type:    "",
        Content: otherResult.Query,
    })
    
    // 合并参数
    for k, v := range otherResult.Parameters {
        q.parameters[k] = v
    }
    
    return q
}

// UnionAll 添加 UNION ALL 子句
func (q *cypherQueryBuilder) UnionAll(other QueryBuilder) QueryBuilder {
    otherResult, _ := other.Build()
    q.clauses = append(q.clauses, clause{
        Type:    UnionAllClause,
        Content: "",
    })
    
    // 添加其他查询的子句
    q.clauses = append(q.clauses, clause{
        Type:    "",
        Content: otherResult.Query,
    })
    
    // 合并参数
    for k, v := range otherResult.Parameters {
        q.parameters[k] = v
    }
    
    return q
}

// SetParameter 设置查询参数
func (q *cypherQueryBuilder) SetParameter(key string, value interface{}) QueryBuilder {
    q.parameters[key] = value
    return q
}

// SetParameters 设置多个查询参数
func (q *cypherQueryBuilder) SetParameters(params map[string]interface{}) QueryBuilder {
    for k, v := range params {
        q.parameters[k] = v
    }
    return q
}

// MatchEntity 匹配实体
func (q *cypherQueryBuilder) MatchEntity(entity interface{}) QueryBuilder {
    pattern, err := q.buildEntityPattern(entity, "")
    if err != nil {
        // 在实际实现中应该返回错误
        return q
    }
    return q.Match(pattern)
}

// CreateEntity 创建实体
func (q *cypherQueryBuilder) CreateEntity(entity interface{}) QueryBuilder {
    pattern, err := q.buildEntityPattern(entity, "")
    if err != nil {
        return q
    }
    return q.Create(pattern)
}

// MergeEntity 合并实体
func (q *cypherQueryBuilder) MergeEntity(entity interface{}) QueryBuilder {
    pattern, err := q.buildEntityPattern(entity, "")
    if err != nil {
        return q
    }
    return q.Merge(pattern)
}

// MatchRelationship 匹配关系
func (q *cypherQueryBuilder) MatchRelationship(from, to interface{}, relType string) QueryBuilder {
    fromPattern, _ := q.buildEntityPattern(from, "from")
    toPattern, _ := q.buildEntityPattern(to, "to")
    
    relBuilder := NewRelationshipBuilder().Type(relType).Direction(DirectionOutgoing)
    relPattern := relBuilder.Build()
    
    fullPattern := fromPattern + relPattern + toPattern
    return q.Match(fullPattern)
}

// CreateRelationship 创建关系
func (q *cypherQueryBuilder) CreateRelationship(from, to interface{}, relType string, props map[string]interface{}) QueryBuilder {
    fromPattern, _ := q.buildEntityPattern(from, "from")
    toPattern, _ := q.buildEntityPattern(to, "to")
    
    relBuilder := NewRelationshipBuilder().Type(relType).Direction(DirectionOutgoing)
    if props != nil {
        relBuilder = relBuilder.Properties(props)
    }
    relPattern := relBuilder.Build()
    
    fullPattern := fromPattern + relPattern + toPattern
    return q.Create(fullPattern)
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
    
    // 构建节点模式
    nodeBuilder := NewNodeBuilder()
    
    if variable != "" {
        nodeBuilder = nodeBuilder.Variable(variable)
    }
    
    nodeBuilder = nodeBuilder.Labels(metadata.Labels...)
    
    // 添加属性
    v := reflect.ValueOf(entity)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    
    properties := make(map[string]interface{})
    for fieldName, propMeta := range metadata.Properties {
        fieldValue := v.FieldByName(fieldName)
        if fieldValue.IsValid() && !fieldValue.IsZero() {
            paramName := q.generateParameterName()
            properties[propMeta.CypherName] = fmt.Sprintf("$%s", paramName)
            q.parameters[paramName] = fieldValue.Interface()
        }
    }
    
    if len(properties) > 0 {
        nodeBuilder = nodeBuilder.Properties(properties)
    }
    
    return nodeBuilder.Build(), nil
}

// generateParameterName 生成参数名
func (q *cypherQueryBuilder) generateParameterName() string {
    q.paramCounter++
    return fmt.Sprintf("param_%d", q.paramCounter)
}

// Build 构建最终的 Cypher 查询
func (q *cypherQueryBuilder) Build() (types.QueryResult, error) {
    var parts []string
    
    for _, clause := range q.clauses {
        part := string(clause.Type)
        if clause.Content != "" {
            if part != "" {
                part += " " + clause.Content
            } else {
                part = clause.Content
            }
        }
        if part != "" {
            parts = append(parts, part)
        }
    }
    
    query := strings.Join(parts, "\n")
    
    return types.QueryResult{
        Query:      query,
        Parameters: q.parameters,
    }, nil
}

// String 返回查询字符串表示
func (q *cypherQueryBuilder) String() string {
    result, _ := q.Build()
    return result.Query
}