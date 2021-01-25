package jsonenv

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Env takes JSON in encoded by Base64 (without padding) and allows read values
// from it.
type Env struct{ root map[string]json.RawMessage }

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

// Read reads string value by path.
func (env Env) Read(fullPath string) (string, error) {
	path := strings.Split(fullPath, "/")
	pathLen := len(path)
	if pathLen == 0 {
		return "", fmt.Errorf(`path "%s" is empty`, fullPath)
	}
	pathLen -= 1

	completedPath := make([]string, 0, pathLen)
	node := env.root
	for i := 0; i < pathLen; i++ {
		key := path[i]
		if len(key) == 0 {
			continue
		}
		child, has := node[key]
		if !has {
			return "", fmt.Errorf(`path node %q is not existing in %q`,
				key, strings.Join(completedPath, "/"))
		}
		completedPath = append(completedPath, key)
		node = nil
		if err := json.Unmarshal(child, &node); err != nil {
			return "", fmt.Errorf(`failed to parse node %q: "%w"`,
				strings.Join(completedPath, "/"), err)
		}
	}

	key := path[pathLen]
	val, has := node[key]
	if !has {
		return "", fmt.Errorf(`path %q doesn't have value key %q`,
			strings.Join(completedPath, "/"), key)
	}

	var result string
	if err := json.Unmarshal(val, &result); err != nil {
		return "", fmt.Errorf(`failed to parse value key %q with path %q: "%w"`,
			key, strings.Join(completedPath, "/"), err)
	}

	return result, nil
}
