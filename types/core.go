// types/core.go
package types

import (
    "time"
)

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

// Path 表示图中的路径
type Path struct {
    Nodes         []Node         `json:"nodes"`
    Relationships []Relationship `json:"relationships"`
    Length        int           `json:"length"`
}

// QueryResult 表示查询构建结果
type QueryResult struct {
    Query      string                 `json:"query"`
    Parameters map[string]interface{} `json:"parameters"`
}

// Pattern 表示 Cypher 模式
type Pattern struct {
    Nodes         []*NodePattern         `json:"nodes"`
    Relationships []*RelationshipPattern `json:"relationships"`
}

// NodePattern 表示节点模式
type NodePattern struct {
    Variable   string                 `json:"variable,omitempty"`
    Labels     []string              `json:"labels"`
    Properties map[string]interface{} `json:"properties"`
}

// RelationshipPattern 表示关系模式
type RelationshipPattern struct {
    Variable   string                 `json:"variable,omitempty"`
    Type       string                `json:"type"`
    Direction  Direction             `json:"direction"`
    Properties map[string]interface{} `json:"properties"`
    MinLength  int                   `json:"minLength,omitempty"`
    MaxLength  int                   `json:"maxLength,omitempty"`
}

// Direction 表示关系方向
type Direction string

const (
    DirectionOutgoing Direction = ">"
    DirectionIncoming Direction = "<"
    DirectionBoth     Direction = ""
)