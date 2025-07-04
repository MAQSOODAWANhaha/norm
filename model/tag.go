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

// ParseTag parses a struct tag string in the format: "name;key1:value1;key2:value2"
func ParseTag(tag string) TagInfo {
	info := TagInfo{
		Options: make(map[string]string),
	}

	parts := strings.Split(tag, ";")
	if len(parts) == 0 {
		return info
	}

	// The first part is the name, unless it contains a colon
	if !strings.Contains(parts[0], ":") {
		info.Name = strings.TrimSpace(parts[0])
		parts = parts[1:]
	}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			info.Options[key] = value
		} else {
			// Handle boolean flags
			info.Options[part] = "true"
		}
	}

	return info
}
