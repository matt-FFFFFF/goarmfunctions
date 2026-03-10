package armparser

import (
	"context"
	"fmt"
	"strings"
)

// resolveMemberAccess applies .dot and ['bracket'] member access to a value.
// bracketMembers are resolved first (they may contain function calls), then dotMembers.
func resolveMemberAccess(value any, dotMembers []string, bracketMembers []*StringExpression, ctx context.Context, evalCtx EvalContext, registry *FuncRegistry) (any, error) {
	if len(dotMembers) == 0 && len(bracketMembers) == 0 {
		return value, nil
	}

	var err error

	// First evaluate square bracket members
	for _, member := range bracketMembers {
		memberValue, err := member.Evaluate(ctx, evalCtx, registry)
		if err != nil {
			return nil, err
		}
		key, ok := memberValue.(string)
		if !ok {
			return nil, fmt.Errorf("member access key must be a string, got %T", memberValue)
		}
		value, err = accessMember(value, key)
		if err != nil {
			return nil, err
		}
	}

	// Next evaluate dot members
	for _, member := range dotMembers {
		value, err = accessMember(value, member)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

// accessMember retrieves a key from a value.
// For map[string]any: direct key lookup, then case-insensitive fallback.
// Missing member returns (nil, nil). Non-map types return an error.
func accessMember(value any, key string) (any, error) {
	m, ok := value.(map[string]any)
	if !ok {
		if value == nil {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot access member %q on non-object type %T", key, value)
	}

	// Direct lookup
	if v, found := m[key]; found {
		return v, nil
	}

	// Case-insensitive fallback
	lowerKey := strings.ToLower(key)
	for k, v := range m {
		if strings.ToLower(k) == lowerKey {
			return v, nil
		}
	}

	return nil, nil
}
