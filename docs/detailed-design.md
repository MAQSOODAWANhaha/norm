# Cypher ORM 简化详细设计文档

## 目录
1. [核心类型和结构](#核心类型和结构)
2. [查询构建器系统](#查询构建器系统)
3. [模型管理系统](#模型管理系统)
4. [类型转换系统](#类型转换系统)
5. [验证系统](#验证系统)
6. [解析系统](#解析系统)
7. [实现指南](#实现指南)
8. [代码示例](#代码示例)
9. [测试指南](#测试指南)

## 核心类型和结构

### 基础类型定义

```go
// types/core.go
package types

import (
    "time"
)

// QueryResult 表示查询构建结果
type QueryResult struct {
    Query      string                 `json:"query"`      // 生成的 Cypher 查询
    Parameters map[string]interface{} `json:"parameters"` // 查询参数
    Valid      bool                   `json:"valid"`      // 查询是否有效 (第二阶段)
    Errors     []ValidationError      `json:"errors"`     // 验证错误 (第二阶段)
}

// ValidationError 表示验证错误
type ValidationError struct {
    Type        string `json:"type"`        // 错误类型
    Message     string `json:"message"`     // 错误消息
    Position    int    `json:"position"`    // 错误位置
    Suggestion  string `json:"suggestion"`  // 修复建议
}

// Node 表示图中的节点
type Node struct {
    ID         interface{}            `json:"id,omitempty"`
    Labels     []string              `json:"labels"`
    Properties map[string]interface{} `json:"properties"`
}

// Relationship 表示图中的关系
type Relationship struct {
    ID         interface{}            `json:"id,omitempty"`
    Type       string                `json:"type"`
    StartNode  interface{}           `json:"startNode"`
    EndNode    interface{}           `json:"endNode"`
    Properties map[string]interface{} `json:"properties"`
}

// Direction 表示关系方向
type Direction string

const (
    DirectionOutgoing Direction = ">"
    DirectionIncoming Direction = "<"
    DirectionBoth     Direction = ""
)

// ClauseType 表示 Cypher 子句类型
type ClauseType string

const (
    MatchClause         ClauseType = "MATCH"
    OptionalMatchClause ClauseType = "OPTIONAL MATCH"
    CreateClause        ClauseType = "CREATE"
    MergeClause         ClauseType = "MERGE"
    SetClause           ClauseType = "SET"
    DeleteClause        ClauseType = "DELETE"
    RemoveClause        ClauseType = "REMOVE"
    ReturnClause        ClauseType = "RETURN"
    WithClause          ClauseType = "WITH"
    WhereClause         ClauseType = "WHERE"
    OrderByClause       ClauseType = "ORDER BY"
    SkipClause          ClauseType = "SKIP"
    LimitClause         ClauseType = "LIMIT"
    UnwindClause        ClauseType = "UNWIND"
)
```

## 查询构建器系统

### 主查询构建器

```go
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
    
    // 参数操作
    SetParameter(key string, value interface{}) QueryBuilder
    SetParameters(params map[string]interface{}) QueryBuilder
    
    // 构建操作
    Build() (types.QueryResult, error)
    Validate() error  // 第二阶段功能
    String() string
    
    // 实体操作
    MatchEntity(entity interface{}) QueryBuilder
    CreateEntity(entity interface{}) QueryBuilder
    MergeEntity(entity interface{}) QueryBuilder
}

// cypherQueryBuilder 实现 QueryBuilder 接口
type cypherQueryBuilder struct {
    clauses       []Clause
    parameters    map[string]interface{}
    paramCounter  int
    registry      *model.Registry
    validator     Validator  // 第二阶段
}

// Clause 表示单个 Cypher 子句
type Clause struct {
    Type    types.ClauseType
    Content string
}

// NewQueryBuilder 创建新的查询构建器
func NewQueryBuilder(registry *model.Registry) QueryBuilder {
    return &cypherQueryBuilder{
        clauses:      make([]Clause, 0),
        parameters:   make(map[string]interface{}),
        paramCounter: 0,
        registry:     registry,
    }
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
    
    result := types.QueryResult{
        Query:      query,
        Parameters: q.parameters,
        Valid:      true,  // 第一阶段默认为 true
        Errors:     nil,
    }
    
    return result, nil
}

// Validate 验证查询 (第二阶段功能)
func (q *cypherQueryBuilder) Validate() error {
    // 第二阶段实现
    return nil
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
```

### 节点构建器

```go
// builder/node.go
package builder

import (
    "fmt"
    "strings"
)

// NodeBuilder 节点构建器接口
type NodeBuilder interface {
    Variable(name string) NodeBuilder
    Labels(labels ...string) NodeBuilder
    Properties(props map[string]interface{}) NodeBuilder
    Property(key string, value interface{}) NodeBuilder
    Build() string
    Clone() NodeBuilder
}

type nodeBuilder struct {
    variable   string
    labels     []string
    properties map[string]interface{}
}

// NewNodeBuilder 创建新的节点构建器
func NewNodeBuilder() NodeBuilder {
    return &nodeBuilder{
        properties: make(map[string]interface{}),
    }
}

// Build 构建节点模式
func (nb *nodeBuilder) Build() string {
    var parts []string
    
    // 添加变量
    if nb.variable != "" {
        parts = append(parts, nb.variable)
    }
    
    // 添加标签
    if len(nb.labels) > 0 {
        labelStr := ":" + strings.Join(nb.labels, ":")
        parts = append(parts, labelStr)
    }
    
    // 添加属性
    if len(nb.properties) > 0 {
        propParts := make([]string, 0, len(nb.properties))
        for k, v := range nb.properties {
            propParts = append(propParts, fmt.Sprintf("%s: %s", k, nb.formatValue(v)))
        }
        propStr := "{" + strings.Join(propParts, ", ") + "}"
        parts = append(parts, propStr)
    }
    
    return "(" + strings.Join(parts, "") + ")"
}

// formatValue 格式化属性值
func (nb *nodeBuilder) formatValue(value interface{}) string {
    switch v := value.(type) {
    case string:
        if strings.HasPrefix(v, "$") {
            return v // 参数引用
        }
        return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "\\'"))
    case int, int64, float64:
        return fmt.Sprintf("%v", v)
    case bool:
        return fmt.Sprintf("%t", v)
    default:
        return fmt.Sprintf("'%v'", v)
    }
}
```

### 关系构建器

```go
// builder/relationship.go
package builder

import (
    "fmt"
    "strings"
    "norm/types"
)

// RelationshipBuilder 关系构建器接口
type RelationshipBuilder interface {
    Variable(name string) RelationshipBuilder
    Type(relType string) RelationshipBuilder
    Direction(dir types.Direction) RelationshipBuilder
    Properties(props map[string]interface{}) RelationshipBuilder
    Property(key string, value interface{}) RelationshipBuilder
    Length(min, max int) RelationshipBuilder
    Build() string
    Clone() RelationshipBuilder
}

type relationshipBuilder struct {
    variable   string
    relType    string
    direction  types.Direction
    properties map[string]interface{}
    minLength  int
    maxLength  int
}

// NewRelationshipBuilder 创建新的关系构建器
func NewRelationshipBuilder() RelationshipBuilder {
    return &relationshipBuilder{
        properties: make(map[string]interface{}),
        minLength:  -1,
        maxLength:  -1,
    }
}

// Build 构建关系模式
func (rb *relationshipBuilder) Build() string {
    var parts []string
    
    // 添加变量
    if rb.variable != "" {
        parts = append(parts, rb.variable)
    }
    
    // 添加类型
    if rb.relType != "" {
        parts = append(parts, ":"+rb.relType)
    }
    
    // 添加长度约束
    if rb.minLength >= 0 || rb.maxLength >= 0 {
        var lengthStr string
        if rb.minLength >= 0 && rb.maxLength >= 0 {
            lengthStr = fmt.Sprintf("*%d..%d", rb.minLength, rb.maxLength)
        } else if rb.minLength >= 0 {
            lengthStr = fmt.Sprintf("*%d..", rb.minLength)
        } else if rb.maxLength >= 0 {
            lengthStr = fmt.Sprintf("*..%d", rb.maxLength)
        }
        parts = append(parts, lengthStr)
    }
    
    // 添加属性
    if len(rb.properties) > 0 {
        propParts := make([]string, 0, len(rb.properties))
        for k, v := range rb.properties {
            propParts = append(propParts, fmt.Sprintf("%s: %s", k, rb.formatValue(v)))
        }
        propStr := "{" + strings.Join(propParts, ", ") + "}"
        parts = append(parts, propStr)
    }
    
    content := strings.Join(parts, "")
    
    // 根据方向格式化
    switch rb.direction {
    case types.DirectionOutgoing:
        return "-[" + content + "]->"
    case types.DirectionIncoming:
        return "<-[" + content + "]-"
    default:
        return "-[" + content + "]-"
    }
}
```

## 模型管理系统

### 实体注册表

```go
// model/registry.go
package model

import (
    "fmt"
    "reflect"
    "strings"
    "sync"
    "norm/types"
)

// Registry 实体注册表接口
type Registry interface {
    Register(entity interface{}) error
    Get(name string) (*EntityMetadata, bool)
    GetByType(t reflect.Type) (*EntityMetadata, bool)
    List() []*EntityMetadata
    Validate(entity interface{}) error
}

// entityRegistry 实现 Registry 接口
type entityRegistry struct {
    entities map[string]*EntityMetadata
    mutex    sync.RWMutex
}

// EntityMetadata 包含实体的元数据
type EntityMetadata struct {
    Type          reflect.Type
    Name          string
    Labels        []string
    Properties    map[string]*PropertyMetadata
    Relationships map[string]*RelationshipMetadata
}

// PropertyMetadata 包含属性的元数据
type PropertyMetadata struct {
    Name         string        // Go 字段名
    CypherName   string        // Cypher 属性名
    Type         reflect.Type  // Go 类型
    CypherType   string        // Cypher 类型
    Required     bool          // 是否必需
    Index        bool          // 是否索引
    Unique       bool          // 是否唯一
    JsonTag      string        // JSON 标签
    Default      interface{}   // 默认值
}

// RelationshipMetadata 包含关系的元数据
type RelationshipMetadata struct {
    Name      string             // Go 字段名
    Type      string             // 关系类型
    Direction types.Direction    // 关系方向
    Target    reflect.Type       // 目标类型
    Multiple  bool               // 是否为多个关系
}

// NewRegistry 创建新的实体注册表
func NewRegistry() Registry {
    return &entityRegistry{
        entities: make(map[string]*EntityMetadata),
    }
}

// Register 注册实体类型
func (er *entityRegistry) Register(entity interface{}) error {
    er.mutex.Lock()
    defer er.mutex.Unlock()
    
    t := reflect.TypeOf(entity)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    
    if t.Kind() != reflect.Struct {
        return fmt.Errorf("entity must be a struct, got %s", t.Kind())
    }
    
    metadata, err := er.extractMetadata(t)
    if err != nil {
        return fmt.Errorf("failed to extract metadata for %s: %w", t.Name(), err)
    }
    
    er.entities[t.Name()] = metadata
    return nil
}

// extractMetadata 从反射类型中提取元数据
func (er *entityRegistry) extractMetadata(t reflect.Type) (*EntityMetadata, error) {
    metadata := &EntityMetadata{
        Type:          t,
        Name:          t.Name(),
        Labels:        er.extractLabels(t),
        Properties:    make(map[string]*PropertyMetadata),
        Relationships: make(map[string]*RelationshipMetadata),
    }
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        
        // 跳过未导出的字段
        if !field.IsExported() {
            continue
        }
        
        // 处理关系字段
        if relTag, ok := field.Tag.Lookup("relationship"); ok {
            rel, err := er.extractRelationship(field, relTag)
            if err != nil {
                return nil, fmt.Errorf("failed to extract relationship %s: %w", field.Name, err)
            }
            metadata.Relationships[field.Name] = rel
            continue
        }
        
        // 处理属性字段
        prop, err := er.extractProperty(field)
        if err != nil {
            return nil, fmt.Errorf("failed to extract property %s: %w", field.Name, err)
        }
        metadata.Properties[field.Name] = prop
    }
    
    return metadata, nil
}
```

## 类型转换系统

### 类型转换器

```go
// types/converter.go
package types

import (
    "fmt"
    "reflect"
    "time"
)

// Converter 类型转换器接口
type Converter interface {
    ToProperty(value interface{}) (interface{}, error)
    FromProperty(value interface{}) (interface{}, error)
    CypherType() string
    Validate(value interface{}) error
}

// ConverterRegistry 类型转换器注册表
type ConverterRegistry struct {
    converters map[reflect.Type]Converter
}

// NewConverterRegistry 创建新的类型转换器注册表
func NewConverterRegistry() *ConverterRegistry {
    registry := &ConverterRegistry{
        converters: make(map[reflect.Type]Converter),
    }
    
    // 注册默认转换器
    registry.registerDefaultConverters()
    
    return registry
}

// registerDefaultConverters 注册内置类型转换器
func (cr *ConverterRegistry) registerDefaultConverters() {
    cr.converters[reflect.TypeOf("")] = &stringConverter{}
    cr.converters[reflect.TypeOf(0)] = &intConverter{}
    cr.converters[reflect.TypeOf(int64(0))] = &int64Converter{}
    cr.converters[reflect.TypeOf(float64(0))] = &float64Converter{}
    cr.converters[reflect.TypeOf(true)] = &boolConverter{}
    cr.converters[reflect.TypeOf(time.Time{})] = &timeConverter{}
}

// Register 注册类型转换器
func (cr *ConverterRegistry) Register(t reflect.Type, converter Converter) {
    cr.converters[t] = converter
}

// GetConverter 获取类型转换器
func (cr *ConverterRegistry) GetConverter(t reflect.Type) (Converter, error) {
    if converter, ok := cr.converters[t]; ok {
        return converter, nil
    }
    return nil, fmt.Errorf("no converter found for type %s", t)
}

// 基础类型转换器实现
type stringConverter struct{}

func (c *stringConverter) ToProperty(value interface{}) (interface{}, error) {
    return value, nil
}

func (c *stringConverter) FromProperty(value interface{}) (interface{}, error) {
    if str, ok := value.(string); ok {
        return str, nil
    }
    return nil, fmt.Errorf("cannot convert %T to string", value)
}

func (c *stringConverter) CypherType() string {
    return "STRING"
}

func (c *stringConverter) Validate(value interface{}) error {
    _, ok := value.(string)
    if !ok {
        return fmt.Errorf("value must be string, got %T", value)
    }
    return nil
}
```

## 验证系统 (第二阶段)

### 查询验证器

```go
// validator/query.go
package validator

import (
    "fmt"
    "strings"
    "norm/types"
)

// QueryValidator 查询验证器接口
type QueryValidator interface {
    Validate(query string) []types.ValidationError
    ValidateStructure(clauses []string) []types.ValidationError
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
            Message:    "查询不能为空",
            Position:   0,
            Suggestion: "请提供有效的 Cypher 查询",
        })
        return errors
    }
    
    // 检查括号匹配
    if !v.validateBrackets(query) {
        errors = append(errors, types.ValidationError{
            Type:       "bracket_mismatch",
            Message:    "括号不匹配",
            Position:   -1,
            Suggestion: "检查所有括号是否正确配对",
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
            if len(stack) == 0 {
                return false
            }
            if stack[len(stack)-1] != pairs[char] {
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
    validClauses := []string{"MATCH", "CREATE", "MERGE", "RETURN", "WITH", "WHERE"}
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
            Message:    "查询必须包含至少一个有效的 Cypher 子句",
            Position:   0,
            Suggestion: "添加 MATCH、CREATE、MERGE 或其他有效子句",
        })
    }
    
    return errors
}
```

## 解析系统 (第二阶段)

### Cypher 解析器

```go
// parser/cypher.go
package parser

import (
    "strings"
    "norm/types"
)

// CypherParser Cypher 解析器接口
type CypherParser interface {
    Parse(query string) (*ParseResult, error)
    ParseClause(clause string) (*ClauseInfo, error)
    ExtractPatterns(query string) ([]PatternInfo, error)
}

// ParseResult 解析结果
type ParseResult struct {
    Clauses  []ClauseInfo  `json:"clauses"`
    Patterns []PatternInfo `json:"patterns"`
    Variables []string     `json:"variables"`
    Parameters []string    `json:"parameters"`
}

// ClauseInfo 子句信息
type ClauseInfo struct {
    Type    types.ClauseType `json:"type"`
    Content string           `json:"content"`
    Line    int              `json:"line"`
    Column  int              `json:"column"`
}

// PatternInfo 模式信息
type PatternInfo struct {
    Type       string                 `json:"type"`       // "node" or "relationship"
    Variable   string                 `json:"variable"`
    Labels     []string              `json:"labels"`
    Properties map[string]interface{} `json:"properties"`
}

// cypherParser 实现 CypherParser 接口
type cypherParser struct{}

// NewCypherParser 创建新的 Cypher 解析器
func NewCypherParser() CypherParser {
    return &cypherParser{}
}

// Parse 解析完整查询
func (p *cypherParser) Parse(query string) (*ParseResult, error) {
    result := &ParseResult{
        Clauses:    make([]ClauseInfo, 0),
        Patterns:   make([]PatternInfo, 0),
        Variables:  make([]string, 0),
        Parameters: make([]string, 0),
    }
    
    lines := strings.Split(query, "\n")
    
    for i, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }
        
        clause, err := p.ParseClause(line)
        if err != nil {
            return nil, err
        }
        
        clause.Line = i + 1
        result.Clauses = append(result.Clauses, *clause)
    }
    
    return result, nil
}

// ParseClause 解析单个子句
func (p *cypherParser) ParseClause(clause string) (*ClauseInfo, error) {
    clause = strings.TrimSpace(clause)
    
    // 识别子句类型
    upperClause := strings.ToUpper(clause)
    
    var clauseType types.ClauseType
    var content string
    
    switch {
    case strings.HasPrefix(upperClause, "MATCH"):
        clauseType = types.MatchClause
        content = strings.TrimSpace(clause[5:])
    case strings.HasPrefix(upperClause, "CREATE"):
        clauseType = types.CreateClause
        content = strings.TrimSpace(clause[6:])
    case strings.HasPrefix(upperClause, "RETURN"):
        clauseType = types.ReturnClause
        content = strings.TrimSpace(clause[6:])
    case strings.HasPrefix(upperClause, "WHERE"):
        clauseType = types.WhereClause
        content = strings.TrimSpace(clause[5:])
    default:
        clauseType = ""
        content = clause
    }
    
    return &ClauseInfo{
        Type:    clauseType,
        Content: content,
    }, nil
}
```

## 实现指南

### 第一阶段实现步骤

1. **创建简化的目录结构**
```bash
# 删除不需要的目录
rm -rf internal cmd tests/integration tests/unit

# 创建新的简化结构
mkdir -p {builder,model,types,validator,parser,examples,tests}
```

2. **实现核心类型** (`types/core.go`)
   - 定义 QueryResult、ValidationError 等基础类型
   - 定义 Direction、ClauseType 等枚举类型

3. **实现查询构建器** (`builder/`)
   - 实现 QueryBuilder 主接口
   - 实现 NodeBuilder 和 RelationshipBuilder
   - 支持基本的 Cypher 子句构建

4. **实现模型系统** (`model/`)
   - 实现 Registry 接口
   - 支持 struct 标签解析
   - 实现实体元数据管理

5. **实现类型转换** (`types/converter.go`)
   - 实现基本类型转换器
   - 支持 Go 类型到 Cypher 类型映射

### 第二阶段实现步骤

1. **实现验证系统** (`validator/`)
   - 实现查询语法验证
   - 实现结构验证
   - 实现参数验证

2. **实现解析系统** (`parser/`)
   - 实现 Cypher 语法解析
   - 实现模式识别
   - 实现变量和参数提取

## 代码示例

### 基本使用示例

```go
// examples/main.go
package main

import (
    "fmt"
    "log"
    "norm/builder"
    "norm/model"
)

// User 用户实体
type User struct {
    ID       int64  `cypher:"id,omitempty"`
    Username string `cypher:"username,required,unique"`
    Email    string `cypher:"email,required"`
    Active   bool   `cypher:"active"`
}

func main() {
    // 创建注册表
    registry := model.NewRegistry()
    
    // 注册实体
    err := registry.Register(User{})
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建查询构建器
    qb := builder.NewQueryBuilder(registry)
    
    // 构建查询
    user := User{
        Username: "johndoe",
        Email:    "john@example.com",
        Active:   true,
    }
    
    result, err := qb.
        CreateEntity(user).
        Return("u").
        Build()
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("生成的查询:")
    fmt.Println(result.Query)
    fmt.Println("参数:")
    for k, v := range result.Parameters {
        fmt.Printf("  %s: %v\n", k, v)
    }
}
```

## 测试指南

### 单元测试结构

```go
// tests/builder_test.go
package tests

import (
    "testing"
    "norm/builder"
    "norm/model"
)

func TestQueryBuilder_Basic(t *testing.T) {
    registry := model.NewRegistry()
    qb := builder.NewQueryBuilder(registry)
    
    result, err := qb.
        Match("(n:Person {name: $name})").
        Return("n").
        SetParameter("name", "test").
        Build()
    
    if err != nil {
        t.Fatal(err)
    }
    
    expected := "MATCH (n:Person {name: $name})\nRETURN n"
    if result.Query != expected {
        t.Errorf("Expected: %s, Got: %s", expected, result.Query)
    }
}
```

这个简化的设计专注于第一和第二阶段的核心功能，去除了复杂的数据库集成部分，使得项目更加轻量和易于维护。