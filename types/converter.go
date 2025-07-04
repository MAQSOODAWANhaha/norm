// types/converter.go
package types

import (
	"fmt"
	"reflect"
	"time"
)

// Converter 类型转换器接口
type Converter interface {
	ToProperty(value interface{}) (interface{}, error)
	FromProperty(value interface{}) (interface{}, error)
	CypherType() string
	Validate(value interface{}) error
}

// ConverterRegistry 类型转换器注册表
type ConverterRegistry struct {
	converters map[reflect.Type]Converter
}

// NewConverterRegistry 创建新的类型转换器注册表
func NewConverterRegistry() *ConverterRegistry {
	registry := &ConverterRegistry{
		converters: make(map[reflect.Type]Converter),
	}

	// 注册默认转换器
	registry.registerDefaultConverters()

	return registry
}

// registerDefaultConverters 注册内置类型转换器
func (cr *ConverterRegistry) registerDefaultConverters() {
	cr.converters[reflect.TypeOf("")] = &stringConverter{}
	cr.converters[reflect.TypeOf(0)] = &intConverter{}
	cr.converters[reflect.TypeOf(int64(0))] = &int64Converter{}
	cr.converters[reflect.TypeOf(float64(0))] = &float64Converter{}
	cr.converters[reflect.TypeOf(true)] = &boolConverter{}
	cr.converters[reflect.TypeOf(time.Time{})] = &timeConverter{}
}

// Register 注册类型转换器
func (cr *ConverterRegistry) Register(t reflect.Type, converter Converter) {
	cr.converters[t] = converter
}

// GetConverter 获取类型转换器
func (cr *ConverterRegistry) GetConverter(t reflect.Type) (Converter, error) {
	if converter, ok := cr.converters[t]; ok {
		return converter, nil
	}
	return nil, fmt.Errorf("no converter found for type %s", t)
}

// 基础类型转换器实现
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
