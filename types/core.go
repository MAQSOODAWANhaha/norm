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

// Clause 表示单个 Cypher 子句
type Clause struct {
	Type    ClauseType
	Content string
}