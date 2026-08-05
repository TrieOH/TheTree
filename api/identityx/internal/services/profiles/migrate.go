package profiles

import (
	"encoding/json"
	"fmt"
)

// pruneToSchema rewrites a profile document to satisfy the given schema by
// dropping every property the schema forbids: keys not listed under
// "properties" are removed wherever "additionalProperties" is false, and
// nested objects are pruned recursively. Fields the schema still allows are
// left untouched. The result is only meaningful if it re-validates against
// the schema; callers must re-run jsonschema.Validate.
func pruneToSchema(schema, doc json.RawMessage) (json.RawMessage, error) {
	var sch map[string]any
	if err := json.Unmarshal(schema, &sch); err != nil {
		return nil, fmt.Errorf("prune: invalid schema json: %w", err)
	}
	var inst any
	if err := json.Unmarshal(doc, &inst); err != nil {
		return nil, fmt.Errorf("prune: invalid instance json: %w", err)
	}
	pruned, err := pruneNode(sch, inst)
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(pruned)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// pruneNode drops keys from an instance object that its schema node
// forbids, recursing into nested objects and array items. It understands
// the JSON Schema draft 2020-12 shape: {"type": ..., "properties": {...},
// "additionalProperties": bool, "items": {...}}. Nodes without a
// "properties" map are returned untouched.
func pruneNode(sch, inst any) (any, error) {
	schemaNode, _ := sch.(map[string]any)
	obj, ok := inst.(map[string]any)
	if !ok {
		if arr, isArr := inst.([]any); isArr {
			items, _ := schemaNode["items"].(map[string]any)
			for i, el := range arr {
				pruned, err := pruneNode(items, el)
				if err != nil {
					return nil, err
				}
				arr[i] = pruned
			}
			return arr, nil
		}
		return inst, nil
	}

	props, _ := schemaNode["properties"].(map[string]any)
	if props == nil {
		return inst, nil
	}
	strict := false
	if ap, ok := schemaNode["additionalProperties"].(bool); ok && !ap {
		strict = true
	}

	out := make(map[string]any, len(obj))
	for key, val := range obj {
		sub, ok := props[key]
		if !ok {
			if strict {
				continue // drop the forbidden key
			}
			out[key] = val
			continue
		}
		pruned, err := pruneNode(sub, val)
		if err != nil {
			return nil, err
		}
		out[key] = pruned
	}
	return out, nil
}
