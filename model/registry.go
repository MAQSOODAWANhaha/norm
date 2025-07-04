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

// NewRegistry 创建新的实体注册表
func NewRegistry() Registry {
	return &entityRegistry{
		entities: make(map[string]*EntityMetadata),
	}
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
	Name       string       // Go 字段名
	CypherName string       // Cypher 属性名
	Type       reflect.Type // Go 类型
	CypherType string       // Cypher 类型
	Required   bool         // 是否必需
	Index      bool         // 是否索引
	Unique     bool         // 是否唯一
	OmitEmpty  bool         // 是否在为空时忽略
}

// RelationshipMetadata 包含关系的元数据
type RelationshipMetadata struct {
	Name      string            // Go 字段名
	Type      string            // 关系类型
	Direction types.Direction   // 关系方向
	Target    reflect.Type      // 目标类型
	Multiple  bool              // 是否为多个关系
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

// Get 获取实体元数据
func (er *entityRegistry) Get(name string) (*EntityMetadata, bool) {
	er.mutex.RLock()
	defer er.mutex.RUnlock()
	metadata, exists := er.entities[name]
	return metadata, exists
}

// GetByType 根据类型获取实体元数据
func (er *entityRegistry) GetByType(t reflect.Type) (*EntityMetadata, bool) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return er.Get(t.Name())
}

// List 列出所有注册的实体
func (er *entityRegistry) List() []*EntityMetadata {
	er.mutex.RLock()
	defer er.mutex.RUnlock()
	var entities []*EntityMetadata
	for _, metadata := range er.entities {
		entities = append(entities, metadata)
	}
	return entities
}

// Validate 验证实体 (暂未实现)
func (er *entityRegistry) Validate(entity interface{}) error {
	// TODO: 实现实体验证逻辑
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
		if !field.IsExported() {
			continue
		}

		if relTag, ok := field.Tag.Lookup("relationship"); ok {
			rel, err := er.extractRelationship(field, relTag)
			if err != nil {
				return nil, fmt.Errorf("failed to extract relationship %s: %w", field.Name, err)
			}
			metadata.Relationships[field.Name] = rel
			continue
		}

		if _, ok := field.Tag.Lookup("cypher"); ok {
			prop, err := er.extractProperty(field)
			if err != nil {
				return nil, fmt.Errorf("failed to extract property %s: %w", field.Name, err)
			}
			metadata.Properties[field.Name] = prop
		}
	}

	return metadata, nil
}

// extractLabels 从 struct tag 中提取标签
func (er *entityRegistry) extractLabels(t reflect.Type) []string {
	// 查找带有 "label" 选项的 "cypher" 标签
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag, ok := field.Tag.Lookup("cypher"); ok {
			tagInfo := ParseTag(tag)
			if labels, ok := tagInfo.Options["label"]; ok {
				return strings.Split(labels, ",")
			}
		}
	}
	// 如果没有找到，则使用类型名称作为默认标签
	return []string{t.Name()}
}

// extractProperty 提取属性元数据
func (er *entityRegistry) extractProperty(field reflect.StructField) (*PropertyMetadata, error) {
	tag := field.Tag.Get("cypher")
	tagInfo := ParseTag(tag)

	prop := &PropertyMetadata{
		Name:       field.Name,
		Type:       field.Type,
		CypherName: tagInfo.Name,
	}

	if prop.CypherName == "" {
		prop.CypherName = strings.ToLower(field.Name)
	}

	if val, ok := tagInfo.Options["unique"]; ok && val == "true" {
		prop.Unique = true
	}
	if val, ok := tagInfo.Options["index"]; ok && val == "true" {
		prop.Index = true
	}
	if val, ok := tagInfo.Options["omitempty"]; ok && val == "true" {
		prop.OmitEmpty = true
	}
	if val, ok := tagInfo.Options["required"]; ok && val == "true" {
		prop.Required = true
	}

	return prop, nil
}

// extractRelationship 提取关系元数据
func (er *entityRegistry) extractRelationship(field reflect.StructField, tag string) (*RelationshipMetadata, error) {
	tagInfo := ParseTag(tag)
	rel := &RelationshipMetadata{
		Name:      field.Name,
		Type:      tagInfo.Name,
		Direction: types.DirectionOutgoing,
		Target:    er.getTargetType(field.Type),
		Multiple:  er.isSliceType(field.Type),
	}

	if dir, ok := tagInfo.Options["direction"]; ok {
		switch strings.ToLower(dir) {
		case "incoming":
			rel.Direction = types.DirectionIncoming
		case "both":
			rel.Direction = types.DirectionBoth
		}
	}
	return rel, nil
}

func (er *entityRegistry) getTargetType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func (er *entityRegistry) isSliceType(t reflect.Type) bool {
	return t.Kind() == reflect.Slice
}