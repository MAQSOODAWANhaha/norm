// builder/entity.go
package builder

import (
	"fmt"
	"reflect"
	"strings"
)

// EntityInfo 包含从实体中解析出的信息
type EntityInfo struct {
	Labels     []string                   // 节点标签
	Properties map[string]interface{}     // 属性键值对
}

// ParseEntity 直接从结构体实例解析标签和属性信息
// 支持的标签格式:
//   - `cypher:"label:User,Admin"`  // 在任意字段上指定标签
//   - `cypher:"username"`          // 字段映射到属性
//   - `cypher:"email,omitempty"`   // 带选项的属性
func ParseEntity(entity interface{}) (*EntityInfo, error) {
	v := reflect.ValueOf(entity)
	t := reflect.TypeOf(entity)

	// 处理指针类型
	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("entity cannot be nil pointer")
		}
		v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity must be a struct, got %s", t.Kind())
	}

	info := &EntityInfo{
		Properties: make(map[string]interface{}),
	}

	// 提取标签和属性
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 跳过未导出的字段
		if !field.IsExported() {
			continue
		}

		cypherTag := field.Tag.Get("cypher")
		if cypherTag == "" {
			continue
		}

		// 解析标签
		tagInfo := parseTag(cypherTag)

		// 处理标签定义
		if labels := tagInfo.getLabels(); len(labels) > 0 {
			info.Labels = labels
			continue
		}

		// 处理属性字段
		if err := extractProperty(info, field, fieldValue, tagInfo); err != nil {
			return nil, fmt.Errorf("failed to extract property %s: %w", field.Name, err)
		}
	}

	// 如果没有找到标签，使用结构体名称作为默认标签
	if len(info.Labels) == 0 {
		info.Labels = []string{t.Name()}
	}

	return info, nil
}

// tagInfo 表示解析后的标签信息
type tagInfo struct {
	name    string
	options map[string]string
}

// parseTag 解析cypher标签
// 支持格式: "property_name,option1,option2" 或 "label:Label1,Label2"
func parseTag(tag string) tagInfo {
	info := tagInfo{
		options: make(map[string]string),
	}

	parts := strings.Split(tag, ",")
	if len(parts) == 0 {
		return info
	}

	// 第一部分可能是属性名或特殊指令
	first := strings.TrimSpace(parts[0])
	if strings.HasPrefix(first, "label:") {
		// 这是标签定义
		labelStr := strings.TrimPrefix(first, "label:")
		info.options["label"] = labelStr
	} else {
		// 这是属性名
		info.name = first
	}

	// 处理其他选项
	for i := 1; i < len(parts); i++ {
		option := strings.TrimSpace(parts[i])
		if option == "" {
			continue
		}

		if strings.Contains(option, ":") {
			// key:value 格式
			kv := strings.SplitN(option, ":", 2)
			if len(kv) == 2 {
				info.options[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		} else {
			// 布尔选项
			info.options[option] = "true"
		}
	}

	return info
}

// getLabels 从标签信息中提取标签列表
func (ti tagInfo) getLabels() []string {
	labelStr, ok := ti.options["label"]
	if !ok {
		return nil
	}

	var labels []string
	for _, label := range strings.Split(labelStr, ",") {
		label = strings.TrimSpace(label)
		if label != "" {
			labels = append(labels, label)
		}
	}
	return labels
}

// hasOption 检查是否有特定选项
func (ti tagInfo) hasOption(key string) bool {
	_, exists := ti.options[key]
	return exists
}

// getOption 获取选项值
func (ti tagInfo) getOption(key string) string {
	return ti.options[key]
}

// extractProperty 从字段中提取属性
func extractProperty(info *EntityInfo, field reflect.StructField, fieldValue reflect.Value, tagInfo tagInfo) error {
	// 确定属性名
	propName := tagInfo.name
	if propName == "" {
		propName = strings.ToLower(field.Name)
	}

	// 检查omitempty选项
	if tagInfo.hasOption("omitempty") && isZeroValue(fieldValue) {
		return nil
	}

	// 获取字段值
	if !fieldValue.IsValid() {
		return nil
	}

	// 转换值
	value, err := convertValue(fieldValue)
	if err != nil {
		return fmt.Errorf("failed to convert value for field %s: %w", field.Name, err)
	}

	info.Properties[propName] = value
	return nil
}

// isZeroValue 检查值是否为零值
func isZeroValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	case reflect.Struct:
		// 对于结构体，使用反射的IsZero方法
		return v.IsZero()
	default:
		return false
	}
}

// convertValue 转换反射值为合适的接口值
func convertValue(v reflect.Value) (interface{}, error) {
	if !v.IsValid() {
		return nil, nil
	}

	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return v.Float(), nil
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.Ptr:
		if v.IsNil() {
			return nil, nil
		}
		return convertValue(v.Elem())
	case reflect.Interface:
		if v.IsNil() {
			return nil, nil
		}
		return v.Interface(), nil
	case reflect.Slice:
		if v.IsNil() {
			return nil, nil
		}
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			elem, err := convertValue(v.Index(i))
			if err != nil {
				return nil, err
			}
			result[i] = elem
		}
		return result, nil
	case reflect.Map:
		if v.IsNil() {
			return nil, nil
		}
		result := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			val, err := convertValue(v.MapIndex(key))
			if err != nil {
				return nil, err
			}
			result[keyStr] = val
		}
		return result, nil
	default:
		// 对于其他类型，直接返回接口值
		return v.Interface(), nil
	}
}

// BuildEntityPattern 从实体信息构建 Cypher 节点模式
func BuildEntityPattern(info *EntityInfo, variable string) string {
	nodeBuilder := NewNodeBuilder()

	if variable != "" {
		nodeBuilder = nodeBuilder.Variable(variable)
	}

	if len(info.Labels) > 0 {
		nodeBuilder = nodeBuilder.Labels(info.Labels...)
	}

	if len(info.Properties) > 0 {
		nodeBuilder = nodeBuilder.Properties(info.Properties)
	}

	return nodeBuilder.Build()
}