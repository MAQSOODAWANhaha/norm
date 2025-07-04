// parser/cypher.go
package parser

import (
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
	Clauses    []ClauseInfo  `json:"clauses"`
	Patterns   []PatternInfo `json:"patterns"`
	Variables  []string      `json:"variables"`
	Parameters []string      `json:"parameters"`
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
	Type       string                 `json:"type"` // "node" or "relationship"
	Variable   string                 `json:"variable"`
	Labels     []string               `json:"labels"`
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
	// 第二阶段实现
	return nil, nil
}

// ParseClause 解析单个子句
func (p *cypherParser) ParseClause(clause string) (*ClauseInfo, error) {
	// 第二阶段实现
	return nil, nil
}

// ExtractPatterns 提取模式
func (p *cypherParser) ExtractPatterns(query string) ([]PatternInfo, error) {
	// 第二阶段实现
	return nil, nil
}
