package validation

import "fmt"

func CheckDeps(deps map[string]any) error {
	var missing []string

	for name, dep := range deps {
		if dep == nil {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required dependencies: %v", missing)
	}

	return nil
}
