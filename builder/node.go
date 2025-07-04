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

// Variable 设置节点变量
func (nb *nodeBuilder) Variable(name string) NodeBuilder {
    nb.variable = name
    return nb
}

// Labels 添加标签
func (nb *nodeBuilder) Labels(labels ...string) NodeBuilder {
    nb.labels = append(nb.labels, labels...)
    return nb
}

// Properties 设置所有属性
func (nb *nodeBuilder) Properties(props map[string]interface{}) NodeBuilder {
    for k, v := range props {
        nb.properties[k] = v
    }
    return nb
}

// Property 设置单个属性
func (nb *nodeBuilder) Property(key string, value interface{}) NodeBuilder {
    nb.properties[key] = value
    return nb
}

// Clone 克隆节点构建器
func (nb *nodeBuilder) Clone() NodeBuilder {
    clone := &nodeBuilder{
        variable:   nb.variable,
        labels:     make([]string, len(nb.labels)),
        properties: make(map[string]interface{}),
    }
    
    copy(clone.labels, nb.labels)
    for k, v := range nb.properties {
        clone.properties[k] = v
    }
    
    return clone
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