// types/core.go
package types

// QueryResult represents the result of a query build.
type QueryResult struct {
	Query      string                 `json:"query"`
	Parameters map[string]interface{} `json:"parameters"`
	Valid      bool                   `json:"valid"`
	Errors     []ValidationError      `json:"errors"`
}

// ValidationError represents a single validation error.
type ValidationError struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Position   int    `json:"position"`
	Suggestion string `json:"suggestion"`
}

// ClauseType represents the type of a Cypher clause.
type ClauseType string

const (
	MatchClause         ClauseType = "MATCH"
	OptionalMatchClause ClauseType = "OPTIONAL MATCH"
	CreateClause        ClauseType = "CREATE"
	MergeClause         ClauseType = "MERGE"
	WhereClause         ClauseType = "WHERE"
	SetClause           ClauseType = "SET"
	DeleteClause        ClauseType = "DELETE"
	DetachDeleteClause  ClauseType = "DETACH DELETE"
	RemoveClause        ClauseType = "REMOVE"
	ReturnClause        ClauseType = "RETURN"
	WithClause          ClauseType = "WITH"
	OrderByClause       ClauseType = "ORDER BY"
	SkipClause          ClauseType = "SKIP"
	LimitClause         ClauseType = "LIMIT"
	OnCreateClause      ClauseType = "ON CREATE"
	OnMatchClause       ClauseType = "ON MATCH"
	UnwindClause        ClauseType = "UNWIND"
	UnionClause         ClauseType = "UNION"
	UnionAllClause      ClauseType = "UNION ALL"
	UseClause           ClauseType = "USE"
	CallClause          ClauseType = "CALL"
	ForEachClause       ClauseType = "FOREACH"
)

// Clause represents a single clause in a Cypher query.
type Clause struct {
	Type    ClauseType
	Content string
}

// Entity is a struct to hold a struct and its alias, used for Return, With, etc.
type Entity struct {
	Struct interface{}
	Alias  string
}

// Pattern represents a graph pattern.
type Pattern struct {
	StartNode    NodePattern
	Relationship RelationshipPattern
	EndNode      NodePattern
}

// Label represents a single node label.
type Label string

// IsValid checks if the label is valid (non-empty).
func (l Label) IsValid() bool {
	return len(l) > 0
}

// Labels represents a collection of node labels.
type Labels []Label

// Contains checks if the collection contains a specific label.
func (ls Labels) Contains(label Label) bool {
	for _, l := range ls {
		if l == label {
			return true
		}
	}
	return false
}

// Add adds a label to the collection if it's not already present and is valid.
func (ls *Labels) Add(label Label) {
	if label.IsValid() && !ls.Contains(label) {
		*ls = append(*ls, label)
	}
}

// Remove removes a label from the collection.
func (ls *Labels) Remove(label Label) {
	for i, l := range *ls {
		if l == label {
			*ls = append((*ls)[:i], (*ls)[i+1:]...)
			return
		}
	}
}

// ToStrings converts the Labels collection to a slice of strings.
func (ls Labels) ToStrings() []string {
	s := make([]string, len(ls))
	for i, l := range ls {
		s[i] = string(l)
	}
	return s
}

// NodePattern represents a node in a pattern.
type NodePattern struct {
	Variable   string
	Labels     Labels
	Properties map[string]interface{}
}

// RelationshipPattern represents a relationship in a pattern.
type RelationshipPattern struct {
	Variable   string
	Type       string
	Direction  RelationshipDirection
	MinLength  *int
	MaxLength  *int
	Properties map[string]interface{}
}

// RelationshipDirection defines the direction of a relationship.
type RelationshipDirection string

const (
	DirectionOutgoing RelationshipDirection = "->"
	DirectionIncoming RelationshipDirection = "<-"
	DirectionBoth     RelationshipDirection = "--"
)

// Direction is an alias for RelationshipDirection for backward compatibility.
type Direction RelationshipDirection
