// builder/relationship.go
package builder

import (
    "fmt"
    "strings"
)

// RelationshipBuilder 关系构建器接口
type RelationshipBuilder interface {
    Variable(name string) RelationshipBuilder
    Type(relType string) RelationshipBuilder
    Direction(dir Direction) RelationshipBuilder
    Properties(props map[string]interface{}) RelationshipBuilder
    Property(key string, value interface{}) RelationshipBuilder
    Length(min, max int) RelationshipBuilder
    Build() string
    Clone() RelationshipBuilder
}

type relationshipBuilder struct {
    variable   string
    relType    string
    direction  Direction
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

// Variable 设置关系变量
func (rb *relationshipBuilder) Variable(name string) RelationshipBuilder {
    rb.variable = name
    return rb
}

// Type 设置关系类型
func (rb *relationshipBuilder) Type(relType string) RelationshipBuilder {
    rb.relType = relType
    return rb
}

// Direction 设置关系方向
func (rb *relationshipBuilder) Direction(dir Direction) RelationshipBuilder {
    rb.direction = dir
    return rb
}

// Properties 设置所有属性
func (rb *relationshipBuilder) Properties(props map[string]interface{}) RelationshipBuilder {
    for k, v := range props {
        rb.properties[k] = v
    }
    return rb
}

// Property 设置单个属性
func (rb *relationshipBuilder) Property(key string, value interface{}) RelationshipBuilder {
    rb.properties[key] = value
    return rb
}

// Length 设置关系长度约束
func (rb *relationshipBuilder) Length(min, max int) RelationshipBuilder {
    rb.minLength = min
    rb.maxLength = max
    return rb
}

// Clone 克隆关系构建器
func (rb *relationshipBuilder) Clone() RelationshipBuilder {
    clone := &relationshipBuilder{
        variable:   rb.variable,
        relType:    rb.relType,
        direction:  rb.direction,
        properties: make(map[string]interface{}),
        minLength:  rb.minLength,
        maxLength:  rb.maxLength,
    }
    
    for k, v := range rb.properties {
        clone.properties[k] = v
    }
    
    return clone
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
    case DirectionOutgoing:
        return "-[" + content + "]->"
    case DirectionIncoming:
        return "<-[" + content + "]-"
    default:
        return "-[" + content + "]-"
    }
}

// formatValue 格式化属性值
func (rb *relationshipBuilder) formatValue(value interface{}) string {
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