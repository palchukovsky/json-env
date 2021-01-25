package jsonenv

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Env takes JSON in encoded by Base64 (without padding) and allows read values
// from it.
type Env struct{ root interface{} }

// NewEnv creates new instance of Env.
func NewEnv(base64Source string) (Env, error) {
	result := Env{}

	sourceJson, err := base64.RawStdEncoding.DecodeString(
		base64Source)
	if err != nil {
		return result, fmt.Errorf(`failed to decode Base64: "%w"`, err)
	}

	if err := json.Unmarshal([]byte(sourceJson), &result.root); err != nil {
		return result, fmt.Errorf(`failed to read JSON: "%w"`, err)
	}

	return result, nil
}

// Dump returns string with content encoded by Base64 (without padding).
func (env Env) Dump() (string, error) {
	source, err := json.Marshal(env.root)
	if err != nil {
		return "", fmt.Errorf(`failed to save in JSON: "%w"`, err)
	}
	return base64.RawStdEncoding.EncodeToString(source), nil
}

// Read reads string value by path.
func (env Env) Read(fullPath string) (*string, error) {
	path := strings.Split(fullPath, "/")
	pathLen := len(path)
	if pathLen == 0 {
		return nil, fmt.Errorf(`path "%s" is empty`, fullPath)
	}
	pathLen -= 1

	node, isNode := env.root.(map[string]interface{})
	if !isNode {
		return nil, fmt.Errorf(`root is not node`)
	}

	completedPath := make([]string, 0, pathLen)
	for i := 0; i < pathLen; i++ {
		key := strings.TrimSpace(path[i])
		if len(key) == 0 {
			continue
		}
		child, has := node[key]
		if !has {
			return nil, nil
		}
		completedPath = append(completedPath, key)
		var isNode bool
		node, isNode = child.(map[string]interface{})
		if !isNode {
			return nil, fmt.Errorf(`failed to path %q is not node`,
				strings.Join(completedPath, "/"))
		}
	}

	key := strings.TrimSpace(path[pathLen])
	val, has := node[key]
	if !has {
		return nil, nil
	}

	result, isString := val.(string)
	if !isString {
		return nil, fmt.Errorf(`value key %q with path %q is not string`,
			key, strings.Join(completedPath, "/"))
	}

	return &result, nil
}

// Set sets value by path.
func (env *Env) Set(fullPath, value string) error {
	path := strings.Split(fullPath, "/")
	pathLen := len(path)
	if pathLen == 0 {
		return fmt.Errorf(`path "%s" is empty`, fullPath)
	}
	pathLen -= 1

	node, isNode := env.root.(map[string]interface{})
	if !isNode {
		return fmt.Errorf(`root is not node`)
	}

	for i := 0; i < pathLen; i++ {
		key := strings.TrimSpace(path[i])
		if len(key) == 0 {
			continue
		}
		if child, has := node[key]; has {
			if childNode, isNode := child.(map[string]interface{}); isNode {
				node = childNode
				continue
			}
		}
		node[key] = map[string]interface{}{}
		node = node[key].(map[string]interface{})
	}

	node[strings.TrimSpace(path[pathLen])] = value

	return nil
}
