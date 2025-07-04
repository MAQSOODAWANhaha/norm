// model/label.go
package model

import (
    "fmt"
    "reflect"
    "strings"
)

// LabelManager 管理节点标签
type LabelManager struct {
    registry *EntityRegistry
}

// NewLabelManager 创建新的标签管理器
func NewLabelManager(registry *EntityRegistry) *LabelManager {
    return &LabelManager{
        registry: registry,
    }
}

// GetLabels 获取实体的标签
func (lm *LabelManager) GetLabels(entity interface{}) ([]string, error) {
    t := reflect.TypeOf(entity)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    
    metadata, exists := lm.registry.GetByType(t)
    if !exists {
        return nil, fmt.Errorf("entity %s not registered", t.Name())
    }
    
    return metadata.Labels, nil
}

// ValidateLabels 验证标签是否有效
func (lm *LabelManager) ValidateLabels(labels []string) error {
    for _, label := range labels {
        if err := lm.validateLabel(label); err != nil {
            return err
        }
    }
    return nil
}

// validateLabel 验证单个标签
func (lm *LabelManager) validateLabel(label string) error {
    if label == "" {
        return fmt.Errorf("label cannot be empty")
    }
    
    // 标签不能包含特殊字符
    if strings.ContainsAny(label, " \t\n\r:()[]{}") {
        return fmt.Errorf("label '%s' contains invalid characters", label)
    }
    
    return nil
}

// FormatLabels 格式化标签为 Cypher 语法
func (lm *LabelManager) FormatLabels(labels []string) string {
    if len(labels) == 0 {
        return ""
    }
    
    return ":" + strings.Join(labels, ":")
}