// builder/relationship.go
package builder

import (
	"fmt"
	"strings"

	"norm/types"
)

// RelationshipBuilder 关系构建器接口
type RelationshipBuilder interface {
	// 基本关系构建
	Type(relType string) RelationshipBuilder
	Variable(variable string) RelationshipBuilder
	Direction(direction types.RelationshipDirection) RelationshipBuilder
	Properties(properties map[string]interface{}) RelationshipBuilder
	
	// 变长路径
	MinLength(min int) RelationshipBuilder
	MaxLength(max int) RelationshipBuilder
	VarLength(min, max int) RelationshipBuilder
	
	// 构建模式
	Build() types.RelationshipPattern
	String() string
}

// relationshipBuilder 关系构建器实现
type relationshipBuilder struct {
	pattern types.RelationshipPattern
}

// NewRelationshipBuilder 创建新的关系构建器
func NewRelationshipBuilder() RelationshipBuilder {
	return &relationshipBuilder{
		pattern: types.RelationshipPattern{
			Direction: types.DirectionOutgoing, // 默认方向
		},
	}
}

// Type 设置关系类型
func (rb *relationshipBuilder) Type(relType string) RelationshipBuilder {
	rb.pattern.Type = relType
	return rb
}

// Variable 设置关系变量
func (rb *relationshipBuilder) Variable(variable string) RelationshipBuilder {
	rb.pattern.Variable = variable
	return rb
}

// Direction 设置关系方向
func (rb *relationshipBuilder) Direction(direction types.RelationshipDirection) RelationshipBuilder {
	rb.pattern.Direction = direction
	return rb
}

// Properties 设置关系属性
func (rb *relationshipBuilder) Properties(properties map[string]interface{}) RelationshipBuilder {
	rb.pattern.Properties = properties
	return rb
}

// MinLength 设置最小长度
func (rb *relationshipBuilder) MinLength(min int) RelationshipBuilder {
	rb.pattern.MinLength = &min
	return rb
}

// MaxLength 设置最大长度
func (rb *relationshipBuilder) MaxLength(max int) RelationshipBuilder {
	rb.pattern.MaxLength = &max
	return rb
}

// VarLength 设置变长路径范围
func (rb *relationshipBuilder) VarLength(min, max int) RelationshipBuilder {
	rb.pattern.MinLength = &min
	rb.pattern.MaxLength = &max
	return rb
}

// Build 构建关系模式
func (rb *relationshipBuilder) Build() types.RelationshipPattern {
	return rb.pattern
}

// String 生成关系模式字符串
func (rb *relationshipBuilder) String() string {
	var sb strings.Builder
	
	// 开始括号和方向
	switch rb.pattern.Direction {
	case types.DirectionIncoming:
		sb.WriteString("<-")
	case types.DirectionOutgoing:
		sb.WriteString("-")
	case types.DirectionBoth:
		sb.WriteString("-")
	default:
		sb.WriteString("-")
	}
	
	sb.WriteString("[")
	
	// 变量名
	if rb.pattern.Variable != "" {
		sb.WriteString(rb.pattern.Variable)
	}
	
	// 关系类型
	if rb.pattern.Type != "" {
		sb.WriteString(":")
		sb.WriteString(rb.pattern.Type)
	}
	
	// 变长路径
	if rb.pattern.MinLength != nil || rb.pattern.MaxLength != nil {
		sb.WriteString("*")
		if rb.pattern.MinLength != nil {
			sb.WriteString(fmt.Sprintf("%d", *rb.pattern.MinLength))
		}
		if rb.pattern.MaxLength != nil {
			sb.WriteString("..")
			sb.WriteString(fmt.Sprintf("%d", *rb.pattern.MaxLength))
		} else if rb.pattern.MinLength != nil {
			sb.WriteString("..")
		}
	}
	
	// 属性 (简化处理，实际应该参数化)
	if len(rb.pattern.Properties) > 0 {
		sb.WriteString(" {")
		var props []string
		for k, v := range rb.pattern.Properties {
			props = append(props, fmt.Sprintf("%s: %v", k, v))
		}
		sb.WriteString(strings.Join(props, ", "))
		sb.WriteString("}")
	}
	
	sb.WriteString("]")
	
	// 结束方向
	switch rb.pattern.Direction {
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

// PatternBuilder 图模式构建器
type PatternBuilder interface {
	StartNode(pattern types.NodePattern) PatternBuilder
	Relationship(pattern types.RelationshipPattern) PatternBuilder
	EndNode(pattern types.NodePattern) PatternBuilder
	Build() types.Pattern
	String() string
}

// patternBuilder 图模式构建器实现
type patternBuilder struct {
	pattern types.Pattern
}

// NewPatternBuilder 创建新的图模式构建器
func NewPatternBuilder() PatternBuilder {
	return &patternBuilder{}
}

// StartNode 设置起始节点
func (pb *patternBuilder) StartNode(pattern types.NodePattern) PatternBuilder {
	pb.pattern.StartNode = pattern
	return pb
}

// Relationship 设置关系
func (pb *patternBuilder) Relationship(pattern types.RelationshipPattern) PatternBuilder {
	pb.pattern.Relationship = pattern
	return pb
}

// EndNode 设置结束节点
func (pb *patternBuilder) EndNode(pattern types.NodePattern) PatternBuilder {
	pb.pattern.EndNode = pattern
	return pb
}

// Build 构建完整模式
func (pb *patternBuilder) Build() types.Pattern {
	return pb.pattern
}

// String 生���完整模式字符串
func (pb *patternBuilder) String() string {
	var sb strings.Builder
	
	// 起始节点
	sb.WriteString(pb.buildNodeString(pb.pattern.StartNode))
	
	// 关系
	rb := &relationshipBuilder{pattern: pb.pattern.Relationship}
	sb.WriteString(rb.String())
	
	// 结束节点
	sb.WriteString(pb.buildNodeString(pb.pattern.EndNode))
	
	return sb.String()
}

// buildNodeString 构建节点字符串
func (pb *patternBuilder) buildNodeString(node types.NodePattern) string {
	var sb strings.Builder
	sb.WriteString("(")
	
	// 变量名
	if node.Variable != "" {
		sb.WriteString(node.Variable)
	}
	
	// 标签
	for _, label := range node.Labels {
		sb.WriteString(":")
		sb.WriteString(label)
	}
	
	// 属性 (简化处理)
	if len(node.Properties) > 0 {
		sb.WriteString(" {")
		var props []string
		for k, v := range node.Properties {
			props = append(props, fmt.Sprintf("%s: %v", k, v))
		}
		sb.WriteString(strings.Join(props, ", "))
		sb.WriteString("}")
	}
	
	sb.WriteString(")")
	return sb.String()
}

// 便利函数

// Outgoing 创建外向关系
func Outgoing(relType string) RelationshipBuilder {
	return NewRelationshipBuilder().Type(relType).Direction(types.DirectionOutgoing)
}

// Incoming 创建入向关系
func Incoming(relType string) RelationshipBuilder {
	return NewRelationshipBuilder().Type(relType).Direction(types.DirectionIncoming)
}

// Bidirectional 创建双向关系
func Bidirectional(relType string) RelationshipBuilder {
	return NewRelationshipBuilder().Type(relType).Direction(types.DirectionBoth)
}

// VarLengthOutgoing 创建变长外向关系
func VarLengthOutgoing(relType string, min, max int) RelationshipBuilder {
	return NewRelationshipBuilder().Type(relType).Direction(types.DirectionOutgoing).VarLength(min, max)
}

// VarLengthIncoming 创建变长入向关系
func VarLengthIncoming(relType string, min, max int) RelationshipBuilder {
	return NewRelationshipBuilder().Type(relType).Direction(types.DirectionIncoming).VarLength(min, max)
}

// VarLengthBidirectional 创建变长双向关系
func VarLengthBidirectional(relType string, min, max int) RelationshipBuilder {
	return NewRelationshipBuilder().Type(relType).Direction(types.DirectionBoth).VarLength(min, max)
}

// Node 创建节点模式
func Node(variable string, labels ...string) types.NodePattern {
	return types.NodePattern{
		Variable: variable,
		Labels:   labels,
	}
}

// NodeWithProps 创建带属性的节点模式
func NodeWithProps(variable string, labels []string, properties map[string]interface{}) types.NodePattern {
	return types.NodePattern{
		Variable:   variable,
		Labels:     labels,
		Properties: properties,
	}
}
