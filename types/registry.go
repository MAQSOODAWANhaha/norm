// types/registry.go
package types

import (
    "fmt"
    "reflect"
    "time"
)

// TypeRegistry 管理类型转换
type TypeRegistry struct {
    converters map[reflect.Type]TypeConverter
}

// TypeConverter 在 Go 和 Cypher 类型之间转换
type TypeConverter interface {
    ToProperty(value interface{}) (interface{}, error)
    FromProperty(value interface{}) (interface{}, error)
    CypherType() string
    Validate(value interface{}) error
}

// NewTypeRegistry 创建新的类型注册表
func NewTypeRegistry() *TypeRegistry {
    registry := &TypeRegistry{
        converters: make(map[reflect.Type]TypeConverter),
    }
    
    // 注册默认转换器
    registry.registerDefaultConverters()
    
    return registry
}

// registerDefaultConverters 注册内置类型转换器
func (tr *TypeRegistry) registerDefaultConverters() {
    tr.converters[reflect.TypeOf("")] = &stringConverter{}
    tr.converters[reflect.TypeOf(0)] = &intConverter{}
    tr.converters[reflect.TypeOf(int64(0))] = &int64Converter{}
    tr.converters[reflect.TypeOf(float64(0))] = &float64Converter{}
    tr.converters[reflect.TypeOf(float32(0))] = &float32Converter{}
    tr.converters[reflect.TypeOf(true)] = &boolConverter{}
    tr.converters[reflect.TypeOf(time.Time{})] = &timeConverter{}
    tr.converters[reflect.TypeOf([]interface{}{})] = &sliceConverter{}
    tr.converters[reflect.TypeOf(map[string]interface{}{})] = &mapConverter{}
}

// Register 注册类型转换器
func (tr *TypeRegistry) Register(t reflect.Type, converter TypeConverter) {
    tr.converters[t] = converter
}

// GetConverter 获取类型转换器
func (tr *TypeRegistry) GetConverter(t reflect.Type) (TypeConverter, error) {
    if converter, ok := tr.converters[t]; ok {
        return converter, nil
    }
    return nil, fmt.Errorf("no converter found for type %s", t)
}

// Convert 转换值为属性值
func (tr *TypeRegistry) Convert(value interface{}) (interface{}, error) {
    t := reflect.TypeOf(value)
    converter, err := tr.GetConverter(t)
    if err != nil {
        return nil, err
    }
    return converter.ToProperty(value)
}

// Validate 验证值是否有效
func (tr *TypeRegistry) Validate(value interface{}) error {
    t := reflect.TypeOf(value)
    converter, err := tr.GetConverter(t)
    if err != nil {
        return err
    }
    return converter.Validate(value)
}

// GetCypherType 获取 Cypher 类型
func (tr *TypeRegistry) GetCypherType(t reflect.Type) string {
    if converter, ok := tr.converters[t]; ok {
        return converter.CypherType()
    }
    return "ANY"
}

// 类型转换器实现

type stringConverter struct{}

func (c *stringConverter) ToProperty(value interface{}) (interface{}, error) {
    return value, nil
}

func (c *stringConverter) FromProperty(value interface{}) (interface{}, error) {
    if str, ok := value.(string); ok {
        return str, nil
    }
    return nil, fmt.Errorf("cannot convert %T to string", value)
}

func (c *stringConverter) CypherType() string {
    return "STRING"
}

func (c *stringConverter) Validate(value interface{}) error {
    _, ok := value.(string)
    if !ok {
        return fmt.Errorf("value must be string, got %T", value)
    }
    return nil
}

type intConverter struct{}

func (c *intConverter) ToProperty(value interface{}) (interface{}, error) {
    return int64(value.(int)), nil
}

func (c *intConverter) FromProperty(value interface{}) (interface{}, error) {
    if i, ok := value.(int64); ok {
        return int(i), nil
    }
    return nil, fmt.Errorf("cannot convert %T to int", value)
}

func (c *intConverter) CypherType() string {
    return "INTEGER"
}

func (c *intConverter) Validate(value interface{}) error {
    _, ok := value.(int)
    if !ok {
        return fmt.Errorf("value must be int, got %T", value)
    }
    return nil
}

type int64Converter struct{}

func (c *int64Converter) ToProperty(value interface{}) (interface{}, error) {
    return value, nil
}

func (c *int64Converter) FromProperty(value interface{}) (interface{}, error) {
    if i, ok := value.(int64); ok {
        return i, nil
    }
    return nil, fmt.Errorf("cannot convert %T to int64", value)
}

func (c *int64Converter) CypherType() string {
    return "INTEGER"
}

func (c *int64Converter) Validate(value interface{}) error {
    _, ok := value.(int64)
    if !ok {
        return fmt.Errorf("value must be int64, got %T", value)
    }
    return nil
}

type float64Converter struct{}

func (c *float64Converter) ToProperty(value interface{}) (interface{}, error) {
    return value, nil
}

func (c *float64Converter) FromProperty(value interface{}) (interface{}, error) {
    if f, ok := value.(float64); ok {
        return f, nil
    }
    return nil, fmt.Errorf("cannot convert %T to float64", value)
}

func (c *float64Converter) CypherType() string {
    return "FLOAT"
}

func (c *float64Converter) Validate(value interface{}) error {
    _, ok := value.(float64)
    if !ok {
        return fmt.Errorf("value must be float64, got %T", value)
    }
    return nil
}

type float32Converter struct{}

func (c *float32Converter) ToProperty(value interface{}) (interface{}, error) {
    return float64(value.(float32)), nil
}

func (c *float32Converter) FromProperty(value interface{}) (interface{}, error) {
    if f, ok := value.(float64); ok {
        return float32(f), nil
    }
    return nil, fmt.Errorf("cannot convert %T to float32", value)
}

func (c *float32Converter) CypherType() string {
    return "FLOAT"
}

func (c *float32Converter) Validate(value interface{}) error {
    _, ok := value.(float32)
    if !ok {
        return fmt.Errorf("value must be float32, got %T", value)
    }
    return nil
}

type boolConverter struct{}

func (c *boolConverter) ToProperty(value interface{}) (interface{}, error) {
    return value, nil
}

func (c *boolConverter) FromProperty(value interface{}) (interface{}, error) {
    if b, ok := value.(bool); ok {
        return b, nil
    }
    return nil, fmt.Errorf("cannot convert %T to bool", value)
}

func (c *boolConverter) CypherType() string {
    return "BOOLEAN"
}

func (c *boolConverter) Validate(value interface{}) error {
    _, ok := value.(bool)
    if !ok {
        return fmt.Errorf("value must be bool, got %T", value)
    }
    return nil
}

type timeConverter struct{}

func (c *timeConverter) ToProperty(value interface{}) (interface{}, error) {
    if t, ok := value.(time.Time); ok {
        return t.Format(time.RFC3339), nil
    }
    return nil, fmt.Errorf("cannot convert %T to time property", value)
}

func (c *timeConverter) FromProperty(value interface{}) (interface{}, error) {
    if str, ok := value.(string); ok {
        return time.Parse(time.RFC3339, str)
    }
    return nil, fmt.Errorf("cannot convert %T to time", value)
}

func (c *timeConverter) CypherType() string {
    return "DATETIME"
}

func (c *timeConverter) Validate(value interface{}) error {
    _, ok := value.(time.Time)
    if !ok {
        return fmt.Errorf("value must be time.Time, got %T", value)
    }
    return nil
}

type sliceConverter struct{}

func (c *sliceConverter) ToProperty(value interface{}) (interface{}, error) {
    return value, nil
}

func (c *sliceConverter) FromProperty(value interface{}) (interface{}, error) {
    if slice, ok := value.([]interface{}); ok {
        return slice, nil
    }
    return nil, fmt.Errorf("cannot convert %T to slice", value)
}

func (c *sliceConverter) CypherType() string {
    return "LIST"
}

func (c *sliceConverter) Validate(value interface{}) error {
    v := reflect.ValueOf(value)
    if v.Kind() != reflect.Slice {
        return fmt.Errorf("value must be slice, got %T", value)
    }
    return nil
}

type mapConverter struct{}

func (c *mapConverter) ToProperty(value interface{}) (interface{}, error) {
    return value, nil
}

func (c *mapConverter) FromProperty(value interface{}) (interface{}, error) {
    if m, ok := value.(map[string]interface{}); ok {
        return m, nil
    }
    return nil, fmt.Errorf("cannot convert %T to map", value)
}

func (c *mapConverter) CypherType() string {
    return "MAP"
}

func (c *mapConverter) Validate(value interface{}) error {
    v := reflect.ValueOf(value)
    if v.Kind() != reflect.Map {
        return fmt.Errorf("value must be map, got %T", value)
    }
    return nil
}