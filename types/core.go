// types/core.go
package types

// QueryResult 表示查询构建结果
type QueryResult struct {
	Query      string                 `json:"query"`      // 生成的 Cypher 查询
	Parameters map[string]interface{} `json:"parameters"` // 查询参数
	Valid      bool                   `json:"valid"`      // 查询是否有效 (第二阶段)
	Errors     []ValidationError      `json:"errors"`     // 验证错误 (第二阶段)
}

// ValidationError 表示验证错误
type ValidationError struct {
	Type       string `json:"type"`       // 错误类型
	Message    string `json:"message"`    // 错误消息
	Position   int    `json:"position"`   // 错误位置
	Suggestion string `json:"suggestion"` // 修复建议
}

// Node 表示图中的节点
type Node struct {
	ID         interface{}            `json:"id,omitempty"`
	Labels     []string               `json:"labels"`
	Properties map[string]interface{} `json:"properties"`
}

// Relationship 表示图中的关系
type Relationship struct {
	ID         interface{} `json:"id,omitempty"`
	Type       string      `json:"type"`
	StartNode  interface{} `json:"startNode"`
	EndNode    interface{} `json:"endNode"`
	Properties map[string]interface{} `json:"properties"`
}

// Path 表示图中的路径
type Path struct {
	Nodes         []Node         `json:"nodes"`
	Relationships []Relationship `json:"relationships"`
	Length        int            `json:"length"`
}

// Direction 表示关系方向
type Direction string

const (
	DirectionOutgoing Direction = ">"
	DirectionIncoming Direction = "<"
	DirectionBoth     Direction = ""
)

// RelationshipPattern 表示关系模式
type RelationshipPattern struct {
	Type       string                 `json:"type"`                 // 关系类型
	Variable   string                 `json:"variable,omitempty"`   // 关系变量名
	Direction  Direction              `json:"direction"`            // 关系方向
	Properties map[string]interface{} `json:"properties,omitempty"` // 关系属性
	MinLength  *int                   `json:"minLength,omitempty"`  // 最小长度 (变长路径)
	MaxLength  *int                   `json:"maxLength,omitempty"`  // 最大长度 (变长路径)
}

// NodePattern 表示节点模式
type NodePattern struct {
	Variable   string                 `json:"variable,omitempty"`   // 节点变量名
	Labels     []string               `json:"labels,omitempty"`     // 节点标签
	Properties map[string]interface{} `json:"properties,omitempty"` // 节点属性
}

// Pattern 表示完整的图模式
type Pattern struct {
	StartNode    NodePattern         `json:"startNode"`    // 起始节点
	Relationship RelationshipPattern `json:"relationship"` // 关系
	EndNode      NodePattern         `json:"endNode"`      // 结束节点
}

// PathLength 表示路径长度范围
type PathLength struct {
	Min *int `json:"min,omitempty"` // 最小长度
	Max *int `json:"max,omitempty"` // 最大长度
}

// ClauseType 表示 Cypher 子句类型
type ClauseType string

const (
	MatchClause         ClauseType = "MATCH"
	OptionalMatchClause ClauseType = "OPTIONAL MATCH"
	CreateClause        ClauseType = "CREATE"
	MergeClause         ClauseType = "MERGE"
	SetClause           ClauseType = "SET"
	DeleteClause        ClauseType = "DELETE"
	DetachDeleteClause  ClauseType = "DETACH DELETE"
	RemoveClause        ClauseType = "REMOVE"
	ReturnClause        ClauseType = "RETURN"
	WithClause          ClauseType = "WITH"
	WhereClause         ClauseType = "WHERE"
	OrderByClause       ClauseType = "ORDER BY"
	SkipClause          ClauseType = "SKIP"
	LimitClause         ClauseType = "LIMIT"
	UnwindClause        ClauseType = "UNWIND"
	CallClause          ClauseType = "CALL"
	UseClause           ClauseType = "USE"
	UnionClause         ClauseType = "UNION"
	UnionAllClause      ClauseType = "UNION ALL"
	ForEachClause       ClauseType = "FOREACH"
	OnCreateClause      ClauseType = "ON CREATE"
	OnMatchClause       ClauseType = "ON MATCH"
)

// Clause 表示单个 Cypher 子句
type Clause struct {
	Type    ClauseType
	Content string
}