// builder/entity.go
package builder

import (
	"fmt"
	"reflect"
	"strings"

	"norm/types"
)

// EntityInfo 存储解析后的实体信息
type EntityInfo struct {
	Labels     types.Labels
	Properties map[string]interface{}
}

// ParseEntity 解析实体结构体，提取标签和属性
func ParseEntity(entity interface{}) (*EntityInfo, error) {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity must be a struct or a pointer to a struct")
	}
	typ := val.Type()

	info := &EntityInfo{
		Properties: make(map[string]interface{}),
	}

	// 1. 解析标签
	info.Labels = parseLabels(typ)

	// 2. 解析属性
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if field.Name == "_" {
			continue
		}

		if !fieldVal.CanInterface() {
			continue
		}

		tag := field.Tag.Get("cypher")
		if tag == "" || tag == "-" {
			continue
		}

		parts := strings.Split(tag, ",")
		propName := parts[0]
		if propName == "" {
			propName = strings.ToLower(field.Name)
		}

		isOmitEmpty := false
		for _, part := range parts {
			if part == "omitempty" {
				isOmitEmpty = true
				break
			}
		}

		if isOmitEmpty && isZero(fieldVal) {
			continue
		}

		info.Properties[propName] = fieldVal.Interface()
	}

	return info, nil
}

// ParseEntityForUpdate 解析实体以进行更新操作
func ParseEntityForUpdate(entity interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity must be a struct or a pointer to a struct")
	}
	typ := val.Type()

	props := make(map[string]interface{})
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if field.Name == "_" || !fieldVal.CanInterface() {
			continue
		}

		tag := field.Tag.Get("cypher")
		if tag == "" || tag == "-" {
			continue
		}

		parts := strings.Split(tag, ",")
		propName := parts[0]
		if propName == "" {
			propName = strings.ToLower(field.Name)
		}

		isOmitEmpty := false
		for _, part := range parts {
			if part == "omitempty" {
				isOmitEmpty = true
				break
			}
		}

		if isOmitEmpty && isZero(fieldVal) {
			continue
		}

		props[propName] = fieldVal.Interface()
	}
	return props, nil
}

// ParseEntityForReturn 解析实体以进行返回操作
func ParseEntityForReturn(entity interface{}, alias string) ([]string, error) {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity must be a struct or a pointer to a struct")
	}
	typ := val.Type()

	var props []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == "_" {
			continue
		}

		tag := field.Tag.Get("cypher")
		if tag == "" || tag == "-" {
			continue
		}

		propName := strings.Split(tag, ",")[0]
		if propName == "" {
			propName = strings.ToLower(field.Name)
		}

		if alias != "" {
			props = append(props, fmt.Sprintf("%s.%s", alias, propName))
		} else {
			props = append(props, propName)
		}
	}
	return props, nil
}

// parseLabels 从类型中解析标签，支持多标签和默认标签生成
func parseLabels(typ reflect.Type) types.Labels {
	var labels types.Labels

	if field, ok := typ.FieldByName("_"); ok {
		tag := field.Tag.Get("cypher")
		if strings.HasPrefix(tag, "label:") {
			labelStr := strings.TrimPrefix(tag, "label:")
			parts := strings.Split(labelStr, ",")
			for _, part := range parts {
				l := types.Label(strings.TrimSpace(part))
				if l.IsValid() {
					labels.Add(l)
				}
			}
		}
	}

	// 如果没有从标签中解析到有效标签，则使用结构体名称作为默认标签
	if len(labels) == 0 {
		defaultLabel := types.Label(typ.Name())
		if defaultLabel.IsValid() {
			labels.Add(defaultLabel)
		}
	}

	return labels
}

// isZero 检查一个 reflect.Value 是���为其类型的零值
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan:
		return v.IsNil()
	case reflect.Bool:
		// false is a valid value and should not be omitted.
		return false
	}
	// For other types like struct, compare against the zero value of that type.
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
