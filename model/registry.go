// model/registry.go
package model

import (
    "fmt"
    "reflect"
    "strings"
    "sync"
    "norm/types"
)

// EntityRegistry 管理实体元数据
type EntityRegistry struct {
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
    Indexes       []*IndexMetadata
    Constraints   []*ConstraintMetadata
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
    Validator    string        // 验证规则
}

// RelationshipMetadata 包含关系的元数据
type RelationshipMetadata struct {
    Name      string             // Go 字段名
    Type      string             // 关系类型
    Direction types.Direction    // 关系方向
    Target    reflect.Type       // 目标类型
    Multiple  bool               // 是否为多个关系
    Lazy      bool               // 是否懒加载
    Cascade   []string           // 级联操作
}

// IndexMetadata 包含索引的元数据
type IndexMetadata struct {
    Name       string   // 索引名称
    Properties []string // 索引属性
    Unique     bool     // 是否唯一索引
    Composite  bool     // 是否复合索引
}

// ConstraintMetadata 包含约束的元数据
type ConstraintMetadata struct {
    Name       string   // 约束名称
    Type       string   // 约束类型（UNIQUE, EXISTS, etc.）
    Properties []string // 约束属性
}

// NewEntityRegistry 创建新的实体注册表
func NewEntityRegistry() *EntityRegistry {
    return &EntityRegistry{
        entities: make(map[string]*EntityMetadata),
    }
}

// Register 注册实体类型
func (er *EntityRegistry) Register(entity interface{}) error {
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

// Get 获取实体元数据
func (er *EntityRegistry) Get(name string) (*EntityMetadata, bool) {
    er.mutex.RLock()
    defer er.mutex.RUnlock()
    
    metadata, exists := er.entities[name]
    return metadata, exists
}

// GetByType 根据类型获取实体元数据
func (er *EntityRegistry) GetByType(t reflect.Type) (*EntityMetadata, bool) {
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    return er.Get(t.Name())
}

// List 列出所有注册的实体
func (er *EntityRegistry) List() []*EntityMetadata {
    er.mutex.RLock()
    defer er.mutex.RUnlock()
    
    var entities []*EntityMetadata
    for _, metadata := range er.entities {
        entities = append(entities, metadata)
    }
    return entities
}

// extractMetadata 从反射类型中提取元数据
func (er *EntityRegistry) extractMetadata(t reflect.Type) (*EntityMetadata, error) {
    metadata := &EntityMetadata{
        Type:          t,
        Name:          t.Name(),
        Labels:        er.extractLabels(t),
        Properties:    make(map[string]*PropertyMetadata),
        Relationships: make(map[string]*RelationshipMetadata),
        Indexes:       er.extractIndexes(t),
        Constraints:   er.extractConstraints(t),
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

// extractLabels 提取标签
func (er *EntityRegistry) extractLabels(t reflect.Type) []string {
    var labels []string
    
    // 从类型标签中提取
    if tag, ok := t.Tag().Lookup("cypher"); ok {
        parts := strings.Split(tag, ",")
        for _, part := range parts {
            part = strings.TrimSpace(part)
            if strings.HasPrefix(part, "label:") {
                label := strings.TrimPrefix(part, "label:")
                labels = append(labels, label)
            }
        }
    }
    
    // 如果没有指定标签，使用结构体名称
    if len(labels) == 0 {
        labels = []string{t.Name()}
    }
    
    return labels
}

// extractProperty 提取属性元数据
func (er *EntityRegistry) extractProperty(field reflect.StructField) (*PropertyMetadata, error) {
    prop := &PropertyMetadata{
        Name:       field.Name,
        Type:       field.Type,
        CypherName: er.getCypherName(field),
        JsonTag:    field.Tag.Get("json"),
    }
    
    // 解析 cypher 标签
    if tag, ok := field.Tag.Lookup("cypher"); ok {
        err := er.parseCypherTag(prop, tag)
        if err != nil {
            return nil, err
        }
    }
    
    // 设置 Cypher 类型
    prop.CypherType = er.getCypherType(field.Type)
    
    return prop, nil
}

// extractRelationship 提取关系元数据
func (er *EntityRegistry) extractRelationship(field reflect.StructField, tag string) (*RelationshipMetadata, error) {
    rel := &RelationshipMetadata{
        Name:     field.Name,
        Target:   er.getTargetType(field.Type),
        Multiple: er.isSliceType(field.Type),
    }
    
    // 解析关系标签
    parts := strings.Split(tag, ",")
    if len(parts) > 0 {
        rel.Type = strings.TrimSpace(parts[0])
    }
    
    for i := 1; i < len(parts); i++ {
        part := strings.TrimSpace(parts[i])
        switch part {
        case "incoming":
            rel.Direction = types.DirectionIncoming
        case "outgoing":
            rel.Direction = types.DirectionOutgoing
        case "lazy":
            rel.Lazy = true
        case "cascade":
            rel.Cascade = append(rel.Cascade, "DELETE")
        }
    }
    
    return rel, nil
}

// extractIndexes 提取索引信息
func (er *EntityRegistry) extractIndexes(t reflect.Type) []*IndexMetadata {
    var indexes []*IndexMetadata
    // 这里可以实现索引提取逻辑
    return indexes
}

// extractConstraints 提取约束信息
func (er *EntityRegistry) extractConstraints(t reflect.Type) []*ConstraintMetadata {
    var constraints []*ConstraintMetadata
    // 这里可以实现约束提取逻辑
    return constraints
}

// getCypherName 获取 Cypher 属性名
func (er *EntityRegistry) getCypherName(field reflect.StructField) string {
    if tag, ok := field.Tag.Lookup("cypher"); ok {
        parts := strings.Split(tag, ",")
        if len(parts) > 0 && parts[0] != "" {
            return parts[0]
        }
    }
    return strings.ToLower(field.Name)
}

// parseCypherTag 解析 cypher 标签
func (er *EntityRegistry) parseCypherTag(prop *PropertyMetadata, tag string) error {
    parts := strings.Split(tag, ",")
    
    // 第一部分是属性名
    if len(parts) > 0 && parts[0] != "" {
        prop.CypherName = parts[0]
    }
    
    // 解析其他选项
    for i := 1; i < len(parts); i++ {
        part := strings.TrimSpace(parts[i])
        switch part {
        case "required":
            prop.Required = true
        case "index":
            prop.Index = true
        case "unique":
            prop.Unique = true
        case "omitempty":
            // 处理 omitempty 标志
        }
    }
    
    return nil
}

// getCypherType 获取 Cypher 类型
func (er *EntityRegistry) getCypherType(t reflect.Type) string {
    switch t.Kind() {
    case reflect.String:
        return "STRING"
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return "INTEGER"
    case reflect.Float32, reflect.Float64:
        return "FLOAT"
    case reflect.Bool:
        return "BOOLEAN"
    case reflect.Slice:
        return "LIST"
    case reflect.Map:
        return "MAP"
    default:
        if t.String() == "time.Time" {
            return "DATETIME"
        }
        return "ANY"
    }
}

// getTargetType 获取目标类型
func (er *EntityRegistry) getTargetType(t reflect.Type) reflect.Type {
    if t.Kind() == reflect.Slice {
        return t.Elem()
    }
    return t
}

// isSliceType 判断是否为切片类型
func (er *EntityRegistry) isSliceType(t reflect.Type) bool {
    return t.Kind() == reflect.Slice
}