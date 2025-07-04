# Cypher ORM 简化详细设计文档

## 目录
1. [核心类型和结构](#核心类型和结构)
2. [查询构建器系统](#查询构建器系统)
3. [模型管理系统](#模型管理系统)
4. [类型转换系统](#类型转换系统)
5. [验证系统](#验证系统)
6. [解析系统](#解析系统)
7. [实现指南](#实现指南)
8. [代码示例](#代码示例)
9. [测试指南](#测试指南)

## 核心类型和结构
// ... (内容保持不变)

## 查询构建器系统
// ... (内容保持不变)

## 模型管理系统

### 实体注册表

The model management system will be enhanced to support a more powerful and flexible struct tag parsing mechanism. This allows for more declarative configuration of entities.

### 标签（Tag）解析

A new, dedicated tag parser will be introduced to handle complex tag definitions. The `cypher` tag will support a key-value format.

**Tag Format:**

`cypher:"[<property_name>];[<key1>:<value1>];[<key2>:<value2>]..."`

*   **`property_name`**: (Optional) The first part of the tag, which defines the property name in the database. If omitted, the lowercase field name is used.
*   **`key:value` pairs**: Semicolon-separated pairs for options.

**Supported Tag Options:**

*   `label`: Specifies the node label. Can be a comma-separated list for multiple labels.
*   `unique`: Marks the property as unique (`true` or `false`).
*   `index`: Marks the property for indexing (`true` or `false`).
*   `omitempty`: Omits the field if it has a zero value.

### `model/registry.go`

The `extractMetadata` function will be updated to use the new tag parser.

```go
// model/registry.go

// ...

// extractMetadata 从反射类型中提取元数据
func (er *entityRegistry) extractMetadata(t reflect.Type) (*EntityMetadata, error) {
    metadata := &EntityMetadata{
        Type:          t,
        Name:          t.Name(),
        Properties:    make(map[string]*PropertyMetadata),
        Relationships: make(map[string]*RelationshipMetadata),
    }

    // New label extraction logic
    labels, err := er.extractLabels(t)
    if err != nil {
        return nil, err
    }
    metadata.Labels = labels

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        // ...
        // Updated property extraction logic using the new tag parser
        prop, err := er.extractProperty(field)
        // ...
    }
    return metadata, nil
}

// extractLabels will now parse the 'cypher' tag on the struct itself (if present)
// or on a dedicated field to determine the label.
func (er *entityRegistry) extractLabels(t reflect.Type) ([]string, error) {
    // ... implementation using the new tag parser ...
}
```

### `model/tag.go`

A new file to house the tag parsing logic.

```go
// model/tag.go
package model

import (
    "strings"
)

// TagInfo holds parsed tag data
type TagInfo struct {
    Name    string
    Options map[string]string
}

// ParseTag parses a struct tag string
func ParseTag(tag string) TagInfo {
    // ... implementation ...
}
```

## 类型转换系统
// ... (内容保持不变)

// ... (其余部分保持不变)